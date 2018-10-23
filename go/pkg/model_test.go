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
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	subs = []byte(`
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
}`)

	notifications = []byte(`
{
   "topics" : [
      {
         "topic" : {
            "name" : "mytopic",
            "user" : "foo"
         },
         "events" : []
      }
   ]
}`)

	events = []byte(`
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
}`)
)

func TestUnmarshal(t *testing.T) {
	var rgwEvents RGWEvents
	var rgwSubs RGWSubscriptions
	var rgwNotifications RGWNotifications

	err := json.Unmarshal(events, &rgwEvents)
	assert.Nil(t, err)
	assert.Equal(t, len(rgwEvents.Events), 1)
	assert.Equal(t, rgwEvents.Events[0].Info.Key.Name, "bar")
	assert.Equal(t, rgwEvents.Events[0].EventType, "OBJECT_CREATE")

	err = json.Unmarshal(subs, &rgwSubs)
	assert.Nil(t, err)
	assert.Equal(t, len(rgwSubs.Subscriptions), 1)
	assert.Equal(t, rgwSubs.Subscriptions[0].Topic.Name, "newtopic")

	err = json.Unmarshal(notifications, &rgwNotifications)
	assert.Nil(t, err)
	assert.Equal(t, len(rgwNotifications.Notifications), 1)
	assert.Equal(t, rgwNotifications.Notifications[0].Topic.Name, "mytopic")
}
