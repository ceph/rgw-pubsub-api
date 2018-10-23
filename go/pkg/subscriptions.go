/*
Copyright 2018 The rgw-pubsub-api Authors

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

package rgwpubsub

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
)

/*
# Subscriptions

- Subscribe to a topic
PUT /subscriptions/<name>?topic=<topic>[&push-endpoint=<endpoint>]

- Get subscription
GET /subscriptions/<name>

- Delete subscription to specific topic
DELETE /subscriptions/<name>

- Pull (fetch) events from a subscription
GET /subscriptions/<name>?events[&max-entries=<max>][&marker=<marker>]

- Ack an event (removes it)
POST /subscriptions/<name>?ack&event-id=<id>
*/

// RGWCreateSubscription creates a subscription for a topic
func (rgw *RGWClient) RGWCreateSubscription(sub, topic, endpoint string) error {
	method := "PUT"
	if len(sub) == 0 || len(topic) == 0 {
		return fmt.Errorf("subscription and topic cannot be empty")
	}
	req_url := rgw.endpoint + "/subscriptions/" + sub + "?topic=" + topic
	if len(endpoint) > 0 {
		req_url = req_url + "&push-endpoint=" + endpoint
	}
	_, err := rgw.rgwDoRequestRaw(method, req_url)
	return err
}

// RGWDeleteSubscription deletes a subscription
func (rgw *RGWClient) RGWDeleteSubscription(sub string) error {
	method := "DELETE"
	if len(sub) == 0 {
		return fmt.Errorf("subscription cannot be empty")
	}
	req_url := rgw.endpoint + "/subscriptions/" + sub
	_, err := rgw.rgwDoRequestRaw(method, req_url)
	return err
}

// RGWDeleteEvent deletes an event from a subscription
func (rgw *RGWClient) RGWDeleteEvent(sub, eventId string) error {
	method := "POST"
	if len(sub) == 0 || len(eventId) == 0 {
		return fmt.Errorf("subscription and eventId cannot be empty")
	}
	req_url := rgw.endpoint + "/subscriptions/" + sub + "?ack&event-id=" + eventId
	_, err := rgw.rgwDoRequestRaw(method, req_url)
	return err
}

// RGWGetEvents gets all events from a subscription
func (rgw *RGWClient) RGWGetEvents(sub string, max int, marker string) (*RGWEvents, error) {
	var events RGWEvents
	method := "GET"
	if len(sub) == 0 {
		return nil, fmt.Errorf("subscription cannot be empty")
	}
	req_url := rgw.endpoint + "/subscriptions/" + sub + "?events"
	if max > 0 {
		req_url = req_url + "&max-entries=" + fmt.Sprintf("%d", max)
	}
	if len(marker) > 0 {
		req_url = req_url + "&marker=" + marker
	}
	out, err := rgw.rgwDoRequestRaw(method, req_url)
	if err != nil {
		if err = json.Unmarshal(out, &events); err != nil {
			glog.Warningf("failed to unmarshal events from %s: %v", string(out), err)
			return nil, err
		}
	}
	return &events, nil
}
