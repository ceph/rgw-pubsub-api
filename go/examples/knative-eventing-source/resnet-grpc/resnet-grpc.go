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

	"github.com/ceph/rgw-pubsub-api/go/examples/knative-eventing-source/resnet-grpc/pkg/resnet"
)

const (
	envResnetServingEndpoint = "RESNET_SERVING_ENDPOINT"
	envAccessId              = "S3_ACCESS_KEY_ID"
	envAccessKey             = "S3_SECRET_ACCESS_KEY"
	envEndpoint              = "S3_HOSTNAME"
)

var (
	bucket          = flag.String("bucket", "", "bucket name")
	key             = flag.String("key", "", "key name")
	testOnly        = flag.Bool("testOnly", false, "test image annotation only, don't start http service")
	servingEndpoint = "localhost:8500"
	rgwDownloader   *rgwdownloader.RGWDownloader
)

func getAnnotation(bucket, key string) {
	reader, err := rgwDownloader.Download(bucket, key)
	if err != nil {
		log.Printf("failed to download %s %s: %v", bucket, key, err)
		return
	}

	tp := resnet.Predict(servingEndpoint, reader)
	log.Printf("classes: %v", tp.Int64Val)
}

func handler(ctx context.Context, e *rgwpubsub.RGWEvent) {
	metadata := cloudevents.FromContext(ctx)
	log.Printf("[%s] %s %s. Object: %q  Bucket: %q", metadata.EventTime.Format(time.RFC3339), metadata.ContentType, metadata.Source, e.Info.Key.Name, e.Info.Bucket.Name)
	getAnnotation(e.Info.Bucket.Name, e.Info.Key.Name)
}

func main() {
	accessId := os.Getenv(envAccessId)
	accessKey := os.Getenv(envAccessKey)
	endpoint := os.Getenv(envEndpoint)
	if len(accessId) == 0 || len(accessKey) == 0 || len(endpoint) == 0 {
		log.Fatalf("env %s, %s, or %s not set", envAccessId, envAccessKey, envEndpoint)
	}

	flag.Parse()
	servingAddr := os.Getenv(envResnetServingEndpoint)
	if len(servingAddr) > 0 {
		servingEndpoint = servingAddr
	}

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
