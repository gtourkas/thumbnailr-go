package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)


func lambdaHandler(ctx context.Context, req events.ALBTargetGroupRequest) (res events.ALBTargetGroupResponse, err error) {

	res.Headers = map[string]string{"content-type":  "application/json"}
	res.IsBase64Encoded = false
	res.StatusCode = http.StatusOK
	res.StatusDescription = fmt.Sprintf("%d %s", res.StatusCode, http.StatusText(res.StatusCode))

	return
}

func main() {
	lambda.Start(lambdaHandler)
}