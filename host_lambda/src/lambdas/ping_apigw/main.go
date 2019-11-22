package main

import (
"context"
"github.com/aws/aws-lambda-go/events"
"github.com/aws/aws-lambda-go/lambda"
"net/http"
)


func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {

	res.Headers = map[string]string{"content-type":  "application/json"}
	res.StatusCode = http.StatusOK

	return
}

func main() {
	lambda.Start(lambdaHandler)
}



