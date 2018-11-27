package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ceph/rgw-pubsub-api/go/pkg"
	"github.com/knative/pkg/cloudevents"

	"github.com/ceph/rgw-pubsub-api/go/examples/knative-eventing-source/pkg/rgwdownloader"

	"github.com/ceph/rgw-pubsub-api/go/examples/knative-eventing-source/vision/pkg/googlevision"
)

const (
	envAPIKey    = "GOOGLE_VISION_API_KEY"
	envAccessId  = "S3_ACCESS_KEY_ID"
	envAccessKey = "S3_SECRET_ACCESS_KEY"
	envEndpoint  = "S3_HOSTNAME"
)

var (
	apiKey    = ""
	numAnnInt = flag.Int("num-annotations", 1, "number of annotations")
	bucket    = flag.String("bucket", "", "bucket name")
	key       = flag.String("key", "", "key name")
	testOnly  = flag.Bool("testOnly", false, "test image annotation only, don't start http service")

	rgwDownloader *rgwdownloader.RGWDownloader
)

func getAnnotation(bucket, key string) {
	reader, err := rgwDownloader.Download(bucket, key)
	if err != nil {
		log.Printf("failed to download %s %s: %v", bucket, key, err)
		return
	}

	if annotations := googlevision.AnnotateImage(apiKey, *numAnnInt, reader); len(annotations) > 0 {
		// print all annotations, highest score first
		for i := 0; i < len(annotations); i++ {
			label := annotations[i].Description
			score := annotations[i].Score
			log.Printf("label: %s, Score: %f\n", label, score)
		}
	}
}

func handler(ctx context.Context, e *rgwpubsub.RGWEvent) {
	metadata := cloudevents.FromContext(ctx)
	log.Printf("[%s] %s %s. Object: %q  Bucket: %q", metadata.EventTime.Format(time.RFC3339), metadata.ContentType, metadata.Source, e.Info.Key.Name, e.Info.Bucket.Name)
	getAnnotation(e.Info.Bucket.Name, e.Info.Key.Name)
}

func main() {
	apiKey = os.Getenv(envAPIKey)
	accessId := os.Getenv(envAccessId)
	accessKey := os.Getenv(envAccessKey)
	endpoint := os.Getenv(envEndpoint)
	if len(accessId) == 0 || len(accessKey) == 0 || len(endpoint) == 0 || len(apiKey) == 0 {
		log.Fatalf("env %s, %s, %s, or %s not set", envAccessId, envAccessKey, envEndpoint, envAPIKey)
	}

	flag.Parse()
	var err error
	rgwDownloader, err = rgwdownloader.NewRGWDownload(accessId, accessKey, endpoint, "default")
	if err != nil {
		log.Fatal(err)
	}

	if bucket != nil && len(*bucket) > 0 && key != nil && len(*key) > 0 {
		getAnnotation(*bucket, *key)
	}
	if !*testOnly {
		log.Print("Ready and listening on port 8080")
		log.Fatal(http.ListenAndServe(":8080", cloudevents.Handler(handler)))
	}
}
