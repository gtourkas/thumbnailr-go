package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"thumbnailr/app/get"
	"thumbnailr/host_lambda/shared"
	"thumbnailr/repos_dynamodb"
	"thumbnailr/stores_s3"
)

var appHandler *get.Handler

func init() {
	sess, err := session.NewSession()
	if err != nil {
		log.Printf("cannot create new sdk session: %s", err)
		return
	}


	thumbnailRepo := repos_dynamodb.NewThumbnailRepo(sess)

	thumbnailStore := stores_s3.NewThumbnailStore(sess,"thumbnailr-thumbnailstore")

	appHandler = &get.Handler{
		ThumbnailRepo: thumbnailRepo,
		ThumbnailStore: thumbnailStore,
	}
}

func lambdaHandler(req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {

	q := req.QueryStringParameters

	in := get.Input{
		UserID: req.RequestContext.Identity.User,
		ThumbnailID: q["id"],
	}

	out := appHandler.Handle(in)

	shared.OutputToResp(&out, &res)

	return
}

func main() {
	lambda.Start(shared.WrapMiddleware(lambdaHandler))
}
