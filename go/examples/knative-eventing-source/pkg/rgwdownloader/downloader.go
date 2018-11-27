package rgwdownloader

import (
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type RGWDownloader struct {
	client *s3.S3
}

func NewRGWDownload(accessKey, secret, endpoint, region string) (*RGWDownloader, error) {
	noSSL := false

	pair := strings.Split(endpoint, "://")

	if len(pair) > 1 {
		noSSL = (pair[0] == "http")
	}

	pathStyle := true
	token := ""
	creds := credentials.NewStaticCredentials(accessKey, secret, token)
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

	if err != nil {
		return nil, fmt.Errorf("Unable to create S3 session instance: %v", err)
	}

	s3Client := s3.New(sess)

	return &RGWDownloader{
		client: s3Client,
	}, nil

}

func (r *RGWDownloader) Download(bucket, key string) (io.ReadCloser, error) {
	var body io.ReadCloser
	resp, err := r.client.GetObject(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return body, err
	}
	body = resp.Body
	return body, nil
}
