package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"net/http"
	"strconv"
	"thumbnailr/app/request_creation"
	"thumbnailr/bus_sns"
	"thumbnailr/host_lambda/shared"
	"thumbnailr/repos_dynamodb"
)

var appHandler *request_creation.Handler

func init() {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("cannot create new sdk session: %s", err)
	}

	cmdIssuer, err := bus_sns.NewCommandIssuer(sess)
	if err != nil {
		log.Fatalf("cannot create command issuer: %s", err)
	}

	thumbnailRepo := repos_dynamodb.NewThumbnailRepo(sess)

	quotaRepo := repos_dynamodb.NewQuotaRepo(sess)

	appHandler = request_creation.NewHandler(quotaRepo, thumbnailRepo, cmdIssuer)
}

func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {

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
		UserID: ctx.Value("UserID").(string),
		PhotoID: q["photoID"],
		Format: q["format"],
		Width: width,
		Length: length,
	}

	out := appHandler.Handle(in)

	shared.OutputToApiGatewayResp(&out, &res)

	return
}

func main() {
	lambda.Start(shared.WrapMiddlewareApiGateway(lambdaHandler))
}
