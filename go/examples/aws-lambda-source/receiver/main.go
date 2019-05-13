package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ceph/rgw-pubsub-api/go/examples/knative-eventing-source/pkg/rgwdownloader"
	"github.com/ceph/rgw-pubsub-api/go/examples/knative-eventing-source/vision/pkg/googlevision"
	"github.com/ceph/rgw-pubsub-api/go/pkg"
	"log"
	"os"
)

type MyResponse struct {
	Message string `json:"answer"`
}

const (
	envS3AccessID         = "S3_ACCESS_KEY_ID"
	envS3AccessKey        = "S3_SECRET_ACCESS_KEY"
	envS3Endpoint         = "S3_HOSTNAME"
	envGoogleVisionAPIKey = "GOOGLE_VISION_API_KEY"
)

var (
	s3AccessID         string
	s3AccessKey        string
	s3Endpoint         string
	googleVisionAPIKey string
	region             = "default"
	maxLabels          = 5
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
	s3AccessID = os.Getenv(envS3AccessID)
	s3AccessKey = os.Getenv(envS3AccessKey)
	s3Endpoint = os.Getenv(envS3Endpoint)
	googleVisionAPIKey = os.Getenv(envGoogleVisionAPIKey)

	if len(s3AccessID) == 0 || len(s3AccessKey) == 0 || len(s3Endpoint) == 0 || len(googleVisionAPIKey) == 0 {
		log.Fatalf("env %s, %s, %s or %s not set", envS3AccessID, envS3AccessKey, envS3Endpoint, envGoogleVisionAPIKey)
	}
	lambda.Start(HandleLambdaEvent)
}
