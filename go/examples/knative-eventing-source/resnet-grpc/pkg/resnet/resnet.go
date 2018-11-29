package resnet

import (
	"context"
	"io"
	"io/ioutil"
	"log"

	tf_core_framework "github.com/tensorflow/tensorflow/tensorflow/go/core/framework"
	pb "tensorflow_serving/apis"

	"google.golang.org/grpc"
)

func Predict(servingAddress string, reader io.ReadCloser) *tf_core_framework.TensorProto {
	imageBytes, err := ioutil.ReadAll(reader)

	if err != nil {
		log.Fatal(err)
	}

	request := &pb.PredictRequest{
		ModelSpec: &pb.ModelSpec{
			Name:          "resnet",
			SignatureName: "serving_default",
		},
		Inputs: map[string]*tf_core_framework.TensorProto{
			"image_bytes": {
				Dtype: tf_core_framework.DataType_DT_STRING,
				TensorShape: &tf_core_framework.TensorShapeProto{
					Dim: []*tf_core_framework.TensorShapeProto_Dim{
						{
							Size: int64(1),
						},
					},
				},
				StringVal: [][]byte{imageBytes},
			},
		},
	}

	conn, err := grpc.Dial(servingAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to the grpc server: %v\n", err)
	}
	defer conn.Close()

	client := pb.NewPredictionServiceClient(conn)

	resp, err := client.Predict(context.Background(), request)
	if err != nil {
		log.Fatalln(err)
	}

	tp, ok := resp.Outputs["classes"]

	if !ok {
		log.Fatalln(err)
	}
	return tp
}
