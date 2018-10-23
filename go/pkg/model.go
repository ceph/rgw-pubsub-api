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

/*
Subscriptions:
{
   "topics" : [
      {
         "topic" : {
            "name" : "newtopic",
            "user" : "foo"
         },
         "subs" : []
      }
   ]
}

Notifications:
{
   "topics" : [
      {
         "topic" : {
            "name" : "newtopic",
            "user" : "foo"
         },
         "events" : []
      }
   ]
}

Events:
{
   "events" : [
      {
         "info" : {
            "attrs" : {
               "mtime" : "2018-10-19 14:00:18.725568Z"
            },
            "bucket" : {
               "tenant" : "",
               "bucket_id" : "f4fa015d-2676-4047-85e7-4a6e1246b5a6.7911.1",
               "name" : "buck"
            },
            "key" : {
               "name" : "bar",
               "instance" : ""
            }
         },
         "timestamp" : "2018-10-18 21:07:08.455516Z2018-10-19 14:00:20.974365Z",
         "id" : "1539957620.974365.924a9e95",
         "event" : "OBJECT_CREATE"
      }
   ],
   "is_truncated" : "false",
   "next_marker" : ""
}
*/

type RGWEventAttr struct {
	Mtime string `json:"mtime,omitempty"`
}

type RGWBucketAttr struct {
	Name     string `json:"name"`
	Tenant   string `json:"tenant,omitempty"`
	BucketId string `json:"bucket_id,omitempty"`
}

type RGWObjectKey struct {
	Name     string `json:"name"`
	Instance string `json:"instance,omitempty"`
}

type RGWEventInfo struct {
	Bucket RGWBucketAttr `json:"bucket"`
	Key    RGWObjectKey  `json:"key"`
	Attrs  RGWEventAttr  `json:"attrs,omitempty"`
}

type RGWEvent struct {
	Info      RGWEventInfo `json:"info"`
	Timestamp string       `json:"timestamp,omitempty"`
	Id        string       `json:"id,omitempty"`
	EventType string       `json:"event,omitempty"`
}

type RGWEvents struct {
	Events      []RGWEvent `json:"events"`
	IsTruncated string     `json:"is_truncated,omitempty"`
	NextMarker  string     `json:"next_marker,omitempty"`
}

type RGWTopic struct {
	Name string `json:"name"`
	User string `json:"user,omitempty"`
}

type RGWSubscription struct {
	Topic        RGWTopic `json:"topic"`
	Subscription []string `json:"subs,omitempty"`
}

type RGWSubscriptions struct {
	Subscriptions []RGWSubscription `json:"topics"`
}

type RGWNotification struct {
	Topic  RGWTopic `json:"topic"`
	Events []string `json:"events,omitempty"`
}

type RGWNotifications struct {
	Notifications []RGWNotification `json:"topics"` /*FIXME: change to notifications */
}
