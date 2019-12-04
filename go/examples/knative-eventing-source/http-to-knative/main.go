package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/knative/pkg/cloudevents"

	"github.com/ceph/rgw-pubsub-api/go/pkg"
)

const (
	envAccessID  = "S3_ACCESS_KEY_ID"
	envAccessKey = "S3_SECRET_ACCESS_KEY"
	envEndpoint  = "S3_HOSTNAME"
)

var (
	userName   = flag.String("username", "", "rgw user name")
	zonegroup  = flag.String("zonegroup", "", "rgw zone group")
	subName    = flag.String("subscriptionname", "", "pubsub subscription name (should exist) for scking")
	target     = flag.String("sink", "", "uri to send events to")
	listenPort = flag.String("port", "8080", "listening port")
	rgwClient  *rgwpubsub.RGWClient
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
	if err := postMessage(*target, &e); err == nil {
		log.Printf("Event %s was successfully posted to knative", e.Id)
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
	flag.Parse()

	if target == nil || *target == "" {
		log.Fatalf("No sink target")
	}

	if subName != nil && len(*subName) > 0 {
		if userName == nil || len(*userName) == 0 || zonegroup == nil || len(*zonegroup) == 0 {
			log.Fatalf("Subscription information exist, but user/zonegroup are missing")
		}
		accessID := os.Getenv(envAccessID)
		accessKey := os.Getenv(envAccessKey)
		endpoint := os.Getenv(envEndpoint)
		if len(accessID) == 0 || len(accessKey) == 0 || len(endpoint) == 0 {
			log.Fatalf("Subscription information exist, but env %s, %s, or %s not set", envAccessID, envAccessKey, envEndpoint)
		}
		var err error
		rgwClient, err = rgwpubsub.NewRGWClient(*userName, accessID, accessKey, endpoint, *zonegroup)
		if err != nil {
			log.Fatalf("Failed to create rgw pubsub client: %v", err)
		}
		log.Printf("Events will acked to rgw: %q", endpoint)
	} else {
		log.Println("No subscription is configured, no events will be acked")
	}

	log.Printf("Listening on port: %q", *listenPort)
	log.Printf("Sink is: %q", *target)

	http.HandleFunc("/", postHandler)
	log.Fatal(http.ListenAndServe(":"+*listenPort, nil))
}

// Creates a CloudEvent Context for a pubsub event.
func cloudEventsContext(e *rgwpubsub.RGWEvent) *cloudevents.EventContext {
	return &cloudevents.EventContext{
		CloudEventsVersion: cloudevents.CloudEventsVersion,
		EventType:          "dev.knative.source.rgwpubsub",
		EventID:            e.Id,
		Source:             "rgwpubsub",
		EventTime:          time.Now(), // use rgw event timestamp?
	}
}

func postMessage(target string, e *rgwpubsub.RGWEvent) error {
	ctx := cloudEventsContext(e)

	log.Printf("Posting to %q", target)
	// Explicitly using Binary encoding so that Istio, et. al. can better inspect
	// event metadata.
	req, err := cloudevents.Binary.NewRequest(target, e, *ctx)
	if err != nil {
		log.Printf("Failed to create http request: %s", err)
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to do POST: %v", err)
		return err
	}
	defer resp.Body.Close()
	log.Printf("response Status: %s", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("response Body: %s", string(body))
	return nil
}
