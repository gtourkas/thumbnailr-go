package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"net/http"
	"strconv"
	"thumbnailr/app/request_creation"
	"thumbnailr/host_lambda/shared"
	"thumbnailr/repos_dynamodb"
)

var appHandler *request_creation.Handler

func init() {
	sess, err := session.NewSession()
	if err != nil {
		log.Printf("cannot create new sdk session: %s", err)
		return
	}

	thumbnailRepo := repos_dynamodb.NewThumbnailRepo(sess)

	quotaRepo := repos_dynamodb.NewQuotaRepo(sess)

	appHandler = &request_creation.Handler{
		ThumbnailRepo: thumbnailRepo,
		QuotaRepo: quotaRepo}
}

func lambdaHandler(req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {

	q := req.QueryStringParameters

	width, err := strconv.Atoi(q["width"])
	if err != nil {
		res.StatusCode = http.StatusBadRequest
		res.Body = "missing or the 'width' parameter is not an integer"
		return
	}
	length, err := strconv.Atoi(q["length"])
	if err != nil {
		res.StatusCode = http.StatusBadRequest
		res.Body = "missing or the 'length' parameter is not an integer"
		return
	}

	in := request_creation.Input{
		UserID: req.RequestContext.Identity.User,
		PhotoID: q["photoID"],
		Format: q["format"],
		Width: width,
		Length: length,
	}

	out := appHandler.Handle(in)

	shared.OutputToResp(&out, &res)

	return
}

func main() {
	lambda.Start(shared.WrapMiddleware(lambdaHandler))
}
