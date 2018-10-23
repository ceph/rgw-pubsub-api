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

package main

import (
	"os"

	"github.com/golang/glog"

	"github.com/ceph/rgw-pubsub-api/go/pkg"
)

const (
	envAccessId  = "S3_ACCESS_KEY_ID"
	envAccessKey = "S3_SECRET_ACCESS_KEY"
	envEndpoint  = "S3_HOSTNAME"
	userName     = "rgwtest"
	zonegroup    = ""
	topicName    = "foobar"
)

func main() {
	accessId := os.Getenv(envAccessId)
	accessKey := os.Getenv(envAccessKey)
	endpoint := os.Getenv(envEndpoint)
	if len(accessId) == 0 || len(accessKey) == 0 || len(endpoint) == 0 {
		glog.Fatalf("env %s, %s, or %s not set", envAccessId, envAccessKey, envEndpoint)
	}
	rgwClient, err := rgwpubsub.NewRGWClient(userName, accessId, accessKey, endpoint, zonegroup)
	if err != nil {
		glog.Fatalf("failed to create rgw pubsub client: %v", err)
	}
	glog.Infof("rgw client %+v", *rgwClient)
	// topic: create, get, delete
	err = rgwClient.RGWCreateTopic(topicName)
	if err != nil {
		glog.Fatalf("failed to create topic: %v", err)
	}
	sub, err := rgwClient.RGWGetSubscriptionWithTopic(topicName)
	if err != nil {
		glog.Fatalf("failed to get sub from topic: %v", err)
	}
	glog.Infof("sub %+v", sub)
	err = rgwClient.RGWDeleteTopic(topicName)
	if err != nil {
		glog.Fatalf("failed to delete topic: %v", err)
	}
}
