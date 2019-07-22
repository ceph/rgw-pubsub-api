package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/ceph/rgw-pubsub-api/go/pkg"
)

const (
	envS3AccessID   = "S3_ACCESS_KEY_ID"
	envS3AccessKey  = "S3_SECRET_ACCESS_KEY"
	envS3Endpoint   = "S3_HOSTNAME"
	envAWSAccessID  = "AWS_ACCESS_KEY_ID"
	envAWSAccessKey = "AWS_SECRET_ACCESS_KEY"
	envAWSRegion    = "AWS_DEFAULT_REGION"
)

var (
	userName     = flag.String("username", "", "rgw user name")
	zonegroup    = flag.String("zonegroup", "", "rgw zone group")
	subName      = flag.String("subscriptionname", "", "pubsub subscription name (should exist) for acking")
	listenPort   = flag.String("port", "8080", "listening port")
	rgwClient    *rgwpubsub.RGWClient
	s3AccessID   string
	s3AccessKey  string
	s3Endpoint   string
	awsAccessID  string
	awsAccessKey string
	awsRegion    string
)

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Printf("%s method not allowed", r.Method)
		http.Error(w, "405 Method Not Allowed", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading message body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print(string(body))

	var e rgwpubsub.RGWEvent
	err = json.Unmarshal(body, &e)

	if err != nil {
		log.Printf("Failed to parse JSON notification: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Successfully created event: %v", e)
	if err := invokeLambda(&e); err == nil {
		log.Printf("Event %s was successfully posted to aws/lambda", e.Id)
		// delete event
		if rgwClient != nil {
			err = rgwClient.RGWDeleteEvent(*subName, e.Id)
			if err != nil {
				log.Printf("Failed to delete event %s: %v", e.Id, err)
			} else {
				log.Printf("Event %s was successfully acked in rgw", e.Id)
			}
		}
	} else {
		log.Printf("Failed to post event %s: %v", e.Id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func main() {
	s3AccessID = os.Getenv(envS3AccessID)
	s3AccessKey = os.Getenv(envS3AccessKey)
	s3Endpoint = os.Getenv(envS3Endpoint)

	if len(s3AccessID) == 0 || len(s3AccessKey) == 0 || len(s3Endpoint) == 0 {
		log.Fatalf("env %s, %s, or %s not set", envS3AccessID, envS3AccessKey, envS3Endpoint)
	}

	awsAccessID = os.Getenv(envAWSAccessID)
	awsAccessKey = os.Getenv(envAWSAccessKey)
	awsRegion = os.Getenv(envAWSRegion)

	if len(awsAccessID) == 0 || len(awsAccessKey) == 0 || len(awsRegion) == 0 {
		log.Fatalf("env %s, %s, or %s not set", envAWSAccessID, envAWSAccessKey, envAWSRegion)
	}
	flag.Parse()

	if subName == nil || len(*subName) == 0 {
		log.Printf("No subscription name - events will not be acked")
	} else {
		var err error
		rgwClient, err = rgwpubsub.NewRGWClient(*userName, s3AccessID, s3AccessKey, s3Endpoint, *zonegroup)
		if err != nil {
			log.Fatalf("Failed to create rgw pubsub client: %v", err)
		}

		log.Printf("Events will acked to rgw: %s", s3Endpoint)
	}
	http.HandleFunc("/", postHandler)
	log.Fatal(http.ListenAndServe(":"+*listenPort, nil))
}

type Response struct {
	Message string `json:"answer"`
}

func invokeLambda(e *rgwpubsub.RGWEvent) error {
	payload, err := json.Marshal(e)
	if err != nil {
		log.Printf("Error marshalling request: %q", err.Error())
		return err
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambda.New(sess, &aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessID, awsAccessKey, ""),
	})

	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String("receiver"), Payload: payload})

	if err != nil {
		log.Printf("Error invoking function: %q", err.Error())
		return err
	}

	if result.StatusCode == nil || *result.StatusCode != 200 {
		log.Printf("Invalid status code: %d", result.StatusCode)
		return nil
	}

	var response Response

	err = json.Unmarshal(result.Payload, &response)
	if err != nil {
		log.Printf("Error unmarshalling receiver response: %q", err)
		return err
	}

	log.Print("response:")
	log.Print(response)
	return nil
}
