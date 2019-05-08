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
	"flag"
	"os"

	"github.com/golang/glog"

	"github.com/ceph/rgw-pubsub-api/go/pkg"
)

const (
	envAccessId  = "S3_ACCESS_KEY_ID"
	envAccessKey = "S3_SECRET_ACCESS_KEY"
	envEndpoint  = "S3_HOSTNAME"
)

func init() {
	flag.Set("logtostderr", "true")
}

var (
	userName         = flag.String("username", "rgwtest", "rgw user name")
	zonegroup        = flag.String("zonegroup", "", "rgw zone group")
	topicName        = flag.String("topicname", "mytopic", "pubsub topic name")
	subName          = flag.String("subname", "mysub", "pubsub subscription name")
	bucketName       = flag.String("bucketname", "buck", "existing rgw bucket name")
	cleanup          = flag.Bool("cleanup", false, "clean up after run")
	readonly         = flag.Bool("readonly", false, "read only")
)

func main() {
	flag.Parse()

	glog.Infof("user name %s, topic %s, sub %s, bucket %s to-clean-up %v readonly %v",
		*userName, *topicName, *subName, *bucketName, *cleanup, *readonly)

	accessId := os.Getenv(envAccessId)
	accessKey := os.Getenv(envAccessKey)
	endpoint := os.Getenv(envEndpoint)
	if len(accessId) == 0 || len(accessKey) == 0 || len(endpoint) == 0 {
		glog.Fatalf("env %s, %s, or %s not set", envAccessId, envAccessKey, envEndpoint)
	}
	rgwClient, err := rgwpubsub.NewRGWClient(*userName, accessId, accessKey, endpoint, *zonegroup)
	if err != nil {
		glog.Fatalf("failed to create rgw pubsub client: %v", err)
	}
	if !*readonly {
		glog.Infof("rgw client %+v", *rgwClient)
		// topic: create, get, delete
		err = rgwClient.RGWCreateTopic(*topicName)
		if err != nil {
			glog.Fatalf("failed to create topic: %v", err)
		}
		// notification: associate a bucket with the topic
		err = rgwClient.RGWCreateNotification(*bucketName, *topicName)
		if err != nil {
			glog.Fatalf("failed to create notification: %v", err)
		}
	}
	notif, err := rgwClient.RGWGetNotifications(*bucketName)
	if err != nil {
		glog.Fatalf("failed to get notifications: %v", err)
	}
	glog.Infof("notifications: %+v", notif)
	if !*readonly {
		// create subscription
		err = rgwClient.RGWCreateSubscription(*subName, *topicName, endpoint)
		if err != nil {
			glog.Fatalf("failed to create subscription: %v", err)
		}
	}
	sub, err := rgwClient.RGWGetSubscriptionWithTopic(*topicName)
	if err != nil {
		glog.Fatalf("failed to get sub from topic: %v", err)
	}
	glog.Infof("subscription: %+v", sub)
	events, err := rgwClient.RGWGetEvents(*subName, 0, "")
	if err != nil {
		glog.Fatalf("failed to get events: %v", err)
	}
	glog.Infof("events: %+v", *events)
	if *cleanup {
		err = rgwClient.RGWDeleteSubscription(*subName)
		if err != nil {
			glog.Fatalf("failed to delete subscription: %v", err)
		}
		err = rgwClient.RGWDeleteNotification(*bucketName, *topicName)
		if err != nil {
			glog.Fatalf("failed to delete notification: %v", err)
		}
		err = rgwClient.RGWDeleteTopic(*topicName)
		if err != nil {
			glog.Fatalf("failed to delete topic: %v", err)
		}
	}
}
