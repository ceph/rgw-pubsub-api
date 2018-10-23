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
# Topics

- Create topic
PUT /topics/<name>

- Get topic
GET /topics/<name>

- List topics
GET /topics

- Delete a topic
DELETE /topics/<name>
*/

// RGWCreateTopic creates a topic
func (rgw *RGWClient) RGWCreateTopic(topic string) error {
	method := "PUT"
	if len(topic) == 0 {
		return fmt.Errorf("topic cannot be empty")
	}
	req_url := rgw.endpoint + "/topics/" + topic
	_, err := rgw.rgwDoRequestRaw(method, req_url)
	return err
}

// RGWDeleteTopic deletes a topic
func (rgw *RGWClient) RGWDeleteTopic(topic string) error {
	method := "DELETE"
	if len(topic) == 0 {
		return fmt.Errorf("topic cannot be empty")
	}
	req_url := rgw.endpoint + "/topics/" + topic
	_, err := rgw.rgwDoRequestRaw(method, req_url)
	return err
}

// RGWGetSubscriptionWithTopic gets a RGWSubscription by providing topic name
func (rgw *RGWClient) RGWGetSubscriptionWithTopic(topic string) (*RGWSubscription, error) {
	var sub RGWSubscription
	method := "GET"
	if len(topic) == 0 {
		return nil, fmt.Errorf("topic cannot be empty")
	}
	req_url := rgw.endpoint + "/topics/" + topic
	out, err := rgw.rgwDoRequestRaw(method, req_url)
	if err != nil {
		if err = json.Unmarshal(out, &sub); err != nil {
			glog.Warningf("failed to unmarshal sub from %s: %v", string(out), err)
			return nil, err
		}
	}
	return &sub, nil
}

// RGWGetSubscriptions gets all subscriptions
func (rgw *RGWClient) RGWGetSubscriptions() (*RGWSubscriptions, error) {
	var subs RGWSubscriptions
	method := "GET"
	req_url := rgw.endpoint + "/topics"
	out, err := rgw.rgwDoRequestRaw(method, req_url)
	if err != nil {
		if err = json.Unmarshal(out, &subs); err != nil {
			glog.Warningf("failed to unmarshal subs from %s: %v", string(out), err)
			return nil, err
		}
	}
	return &subs, nil
}
