/*
Copyright 2018 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/knative/pkg/cloudevents"

	"github.com/ceph/rgw-pubsub-api/go/pkg"
)

const (
	envAccessId  = "S3_ACCESS_KEY_ID"
	envAccessKey = "S3_SECRET_ACCESS_KEY"
	envEndpoint  = "S3_HOSTNAME"
)

var (
	userName  = flag.String("username", "", "rgw user name")
	zonegroup = flag.String("zonegroup", "", "rgw zone group")
	subName   = flag.String("subscriptionname", "", "pubsub subscription name")
	target    = flag.String("sink", "", "uri to send events to")
	pollInt   = flag.String("interval", "5", "polling interval in seconds")
)

type PubSubEventWatcher struct {
	client *rgwpubsub.RGWClient
	target string
}

func NewPubSubEventWatcher(client *rgwpubsub.RGWClient, target string) *PubSubEventWatcher {
	return &PubSubEventWatcher{target: target, client: client}
}

func (ps *PubSubEventWatcher) updateEvent(old, new interface{}) {
	ps.addEvent(new)
}

func (ps *PubSubEventWatcher) addEvent(new interface{}) {
	event := new.(*rgwpubsub.RGWEvent)
	log.Printf("GOT EVENT: %+v", event)
	postMessage(ps.target, event)
}

func main() {
	accessId := os.Getenv(envAccessId)
	accessKey := os.Getenv(envAccessKey)
	endpoint := os.Getenv(envEndpoint)
	if len(accessId) == 0 || len(accessKey) == 0 || len(endpoint) == 0 {
		log.Fatalf("env %s, %s, or %s not set", envAccessId, envAccessKey, envEndpoint)
	}

	flag.Parse()
	if subName == nil || len(*subName) == 0 {
		log.Fatalf("No subscription name")
	}

	if target == nil || *target == "" {
		log.Fatalf("No sink target")
	}

	rgwClient, err := rgwpubsub.NewRGWClient(*userName, accessId, accessKey, endpoint, *zonegroup)
	if err != nil {
		glog.Fatalf("failed to create rgw pubsub client: %v", err)
	}

	log.Printf("Target is: %q", *target)

	var period time.Duration
	if p, err := strconv.Atoi(*pollInt); err != nil {
		period = time.Duration(5) * time.Second
	} else {
		period = time.Duration(p) * time.Second
	}

	ticker := time.NewTicker(period)
	for {
		events, err := rgwClient.RGWGetEvents(*subName, 0, "")
		if err != nil {
			log.Printf("failed to gets event: %v", err)
		} else {
			if events != nil {
				for _, e := range events.Events {
					err = postMessage(*target, &e)
					if err == nil {
						// delete event
						err = rgwClient.RGWDeleteEvent(*subName, e.Id)
						if err != nil {
							log.Printf("failed to delete event %s: %v", e.Id, err)
						}
					} else {
						log.Printf("failed to post event %s: %v", e.Id, err)
					}
				}
			}
		}
		// Wait for next tick
		<-ticker.C
	}
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
