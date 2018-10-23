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
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/service/s3"
)

type RGWUser struct {
	name      string
	accessKey string
	secret    string
}

type RGWClient struct {
	endpoint  string
	zonegroup string
	user      RGWUser
	client    *s3.S3
}

func NewRGWClient(userName, accessKey, secret, endpoint, zonegroup string) (*RGWClient, error) {
	user := RGWUser{
		name:      userName,
		accessKey: accessKey,
		secret:    secret,
	}
	s3, err := getS3Client(user, endpoint, zonegroup)
	if err != nil {
		return nil, err
	}
	return &RGWClient{
		endpoint:  endpoint,
		zonegroup: zonegroup,
		user:      user,
		client:    s3,
	}, nil
}

func getS3Client(user RGWUser, endpoint, region string) (*s3.S3, error) {
	glog.V(5).Infof("Creating s3 client based on: \"%s\" on endpoint %s (%s)", user.accessKey, endpoint, region)

	addr := endpoint
	noSSL := false

	pair := strings.Split(endpoint, "://")

	if len(pair) > 1 {
		noSSL = (pair[0] == "http")
		addr = pair[1]
	}

	pathStyle := true
	token := ""
	creds := credentials.NewStaticCredentials(user.accessKey, user.secret, token)
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:           &region,
			Credentials:      creds,
			Endpoint:         &endpoint,
			DisableSSL:       &noSSL,
			S3ForcePathStyle: &pathStyle,
		},
		// Profile: "profile_name",
	})

	glog.V(5).Infof("  addr=%s (ssl=%t)", addr, !noSSL)

	if err != nil {
		return nil, fmt.Errorf("Unable to create S3 session instance: %v", err)
	}

	s3Client := s3.New(sess)

	return s3Client, nil

}

func (rgw *RGWClient) rgwDoRequestRaw(method, req_url string) ([]byte, error) {
	httpClient := &http.Client{
		Timeout:   30 * time.Second,
		Transport: http.DefaultTransport,
	}

	glog.V(5).Infof("sending http request: %s", req_url)

	req, err := http.NewRequest(method, req_url, bytes.NewReader(nil))
	if err != nil {
		glog.Warningf("Error creating http request %v", err)
		return nil, fmt.Errorf("Error creating http request %v", err)
	}

	token := ""
	s := v4.NewSigner(credentials.NewStaticCredentials(rgw.user.accessKey, rgw.user.secret, token))

	_, err = s.Sign(req, nil, "s3", "default", time.Now())
	if req.Header.Get("Authorization") == "" {
		glog.Warningf("Error signing request: Authorization header is missing")
		return nil, fmt.Errorf("Error signing request: Authorization header is missing")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		glog.Warningf("Error sending http request: %v", err)
		return nil, fmt.Errorf("httpClient.Do err: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		glog.Warningf("Error got http resonse: %v", resp.StatusCode)
		return nil, fmt.Errorf("Error got http response: %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Warningf("Error reading response: %v", err)
		return nil, fmt.Errorf("Error reading response: %v", err)
	}

	return body, nil

}
