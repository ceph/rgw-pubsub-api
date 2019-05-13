package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:name`
	Age  int    `json:age`
}

type MyResponse struct {
	Message string `json:"Answer:"`
}

func HandleLambdaEvent(event MyEvent) (MyResponse, error) {
	return MyResponse{Message: fmt.Sprintf("%s is %d years old! yeehaa!!!!!!!!", event.Name, event.Age)}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
