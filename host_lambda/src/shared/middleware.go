package shared

import (
	"github.com/aws/aws-lambda-go/events"
	"thumbnailr/app"
)

type ApiGatewayHandler func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

var auth app.Auth

func init() {
	auth = app.Auth{PrivateKey:"no-key"}
}

func WrapMiddleware(f ApiGatewayHandler) ApiGatewayHandler {
	return AddAuth(&auth,f)
}
