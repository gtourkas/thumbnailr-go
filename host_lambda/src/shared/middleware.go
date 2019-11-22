package shared

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"thumbnailr/app"
)

var auth app.Auth

func init() {
	auth = app.Auth{PrivateKey:"no-key"}
}

type ApiGatewayHandler func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func WrapMiddlewareApiGateway(f ApiGatewayHandler) ApiGatewayHandler {
	return AddAuthApiGateway(&auth,f)
}

type AlbTargetGroupHandler func(context.Context, events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error)

func WrapMiddlewareAlbTargetGroup(f AlbTargetGroupHandler) AlbTargetGroupHandler {
	return AddAuthAlbTargetGroup(&auth,f)
}

