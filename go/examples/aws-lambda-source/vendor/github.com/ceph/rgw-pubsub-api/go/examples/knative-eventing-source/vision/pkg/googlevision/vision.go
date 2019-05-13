package googlevision

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/vision/v1"
)

func AnnotateImage(apiKey string, numAnnInt int, reader io.ReadCloser) []*vision.EntityAnnotation {
	b, err := ioutil.ReadAll(reader)

	if err != nil {
		log.Fatal(err)
	}

	req := &vision.AnnotateImageRequest{
		// base64 encode
		Image: &vision.Image{
			Content: base64.StdEncoding.EncodeToString(b),
		},
		Features: []*vision.Feature{
			{
				Type:       "LABEL_DETECTION",
				MaxResults: int64(numAnnInt),
			},
		},
	}

	batch := &vision.BatchAnnotateImagesRequest{
		Requests: []*vision.AnnotateImageRequest{req},
	}

	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKey},
	}
	service, err := vision.New(client)
	if err != nil {
		log.Fatal(err)
	}

	res, err := service.Images.Annotate(batch).Do()
	if err != nil {
		log.Fatal(err)
	}
	return res.Responses[0].LabelAnnotations
}
