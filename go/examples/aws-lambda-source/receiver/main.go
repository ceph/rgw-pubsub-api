package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ceph/rgw-pubsub-api/go/examples/knative-eventing-source/pkg/rgwdownloader"
	"github.com/ceph/rgw-pubsub-api/go/examples/knative-eventing-source/vision/pkg/googlevision"
	"github.com/ceph/rgw-pubsub-api/go/pkg"
	"log"
)

type MyResponse struct {
	Message string `json:"answer"`
}

//go:generate ./secret_vars.sh
var (
	region    = "default"
	maxLabels = 5
)

func HandleLambdaEvent(event rgwpubsub.RGWEvent) (MyResponse, error) {
	bucket := event.Info.Bucket.Name
	key := event.Info.Key.Name
	if len(bucket) == 0 || len(key) == 0 {
		message := "missing bucket/key values"
		log.Print(message)
		return MyResponse{}, fmt.Errorf(message)
	}

	log.Print("handle event:")
	log.Print(event)

	rgwDownloader, err := rgwdownloader.NewRGWDownload(s3AccessID, s3AccessKey, s3Endpoint, region)
	if err != nil {
		message := fmt.Sprintf("failed to create downloader: %q", err.Error())
		log.Print(message)
		return MyResponse{}, fmt.Errorf(message)
	}

	reader, err := rgwDownloader.Download(bucket, key)
	if err != nil {
		message := fmt.Sprintf("failed to download object: %q", err.Error())
		log.Printf(message)
		return MyResponse{}, fmt.Errorf(message)
	}

	if annotations := googlevision.AnnotateImage(googleVisionAPIKey, maxLabels, reader); len(annotations) > 0 {
		// return highest score descriptions
		message := fmt.Sprintf("%s/%s is classified as: ", bucket, key)
		for i := 0; i < len(annotations); i++ {
			message += annotations[i].Description + ","
		}
		log.Print(message)
		return MyResponse{Message: message}, nil
	}

	message := "failed to get annotations"
	log.Print(message)
	return MyResponse{}, fmt.Errorf("message")
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
