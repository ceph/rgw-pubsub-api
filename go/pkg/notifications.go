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
# Notifications

GET /notifications/bucket/<bucket>
PUT /notifications/bucket/<bucket>?topic=<topic>
DELETE /notifications/bucket/<bucket>?topic=<topic>
*/

// RGWCreateNotification creates a notification for a topic on a bucket
func (rgw *RGWClient) RGWCreateNotification(bucket, topic string) error {
	method := "PUT"
	if len(bucket) == 0 || len(topic) == 0 {
		return fmt.Errorf("bucket and topic cannot be empty")
	}
	// FIXME: is it URL right?
	req_url := rgw.endpoint + "/notifications/bucket/" + bucket + "?topic=" + topic
	_, err := rgw.rgwDoRequestRaw(method, req_url)
	return err
}

// RGWDeleteNotification deletes a notification
func (rgw *RGWClient) RGWDeleteNotification(bucket, topic string) error {
	method := "DELETE"
	if len(bucket) == 0 || len(topic) == 0 {
		return fmt.Errorf("bucket and topic cannot be empty")
	}
	// FIXME: is it URL right?
	req_url := rgw.endpoint + "/notifications/bucket/" + bucket + "?topic=" + topic
	_, err := rgw.rgwDoRequestRaw(method, req_url)
	return err
}

// RGWGetNotifications gets all notifications on a bucket
func (rgw *RGWClient) RGWGetNotifications(bucket string) (*RGWNotifications, error) {
	var topics RGWNotifications
	method := "GET"
	if len(bucket) == 0 {
		return nil, fmt.Errorf("bucket cannot be empty")
	}

	req_url := rgw.endpoint + "/notifications/bucket/" + bucket
	out, err := rgw.rgwDoRequestRaw(method, req_url)
	if err != nil {
		if err = json.Unmarshal(out, &topics); err != nil {
			glog.Warningf("failed to unmarshal topics from %s: %v", string(out), err)
			return nil, err
		}
	}
	return &topics, nil
}
