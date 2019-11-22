package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"thumbnailr/app/check_creation"
	"thumbnailr/host_lambda/shared"
	"thumbnailr/repos_dynamodb"
)

var appHandler *check_creation.Handler

func init() {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("cannot create new sdk session: %s", err)
	}

	thumbnailRepo := repos_dynamodb.NewThumbnailRepo(sess)

	appHandler = &check_creation.Handler{
		ThumbnailRepo: thumbnailRepo,
	}
}

func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {

	q := req.QueryStringParameters

	in := check_creation.Input{
		UserID: ctx.Value("UserID").(string),
		ThumbnailID: q["id"],
	}

	out := appHandler.Handle(in)

	shared.OutputToApiGatewayResp(&out, &res)

	return
}

func main() {
	lambda.Start(shared.WrapMiddlewareApiGateway(lambdaHandler))
}
