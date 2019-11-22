package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
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
		log.Fatalf("cannot create new sdk session: %s", err)
	}

	svc := sts.New(sess)
	cid, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatalf("cannot get the caller identity: %s", err)
	}

	thumbnailRepo := repos_dynamodb.NewThumbnailRepo(sess)

	thumbnailStore := stores_s3.NewThumbnailStore(sess, fmt.Sprintf("thumbnailr-thumbnailstore-%s", *cid.Account))

	appHandler = &get.Handler{
		ThumbnailRepo: thumbnailRepo,
		ThumbnailStore: thumbnailStore,
	}
}

func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {

	q := req.QueryStringParameters

	in := get.Input{
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
