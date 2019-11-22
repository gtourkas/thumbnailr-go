package shared

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"net/http"
	"thumbnailr/app"
)

func AddAuthApiGateway(auth *app.Auth, h ApiGatewayHandler) ApiGatewayHandler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {

		authHeader := req.Headers["Authorization"]

		userID, err := auth.ParseAuthHeader(authHeader)
		if err != nil {
			res.StatusCode = http.StatusInternalServerError
			e := errors.Wrap(err, "cannot parse the access token")
			res.Body = e.Error()
			return
		}

		ctx = context.WithValue(ctx, "UserID", userID)

		res, err = h(ctx, req)
		return
	}
}

func AddAuthAlbTargetGroup(auth *app.Auth, h AlbTargetGroupHandler) AlbTargetGroupHandler {
	return func(ctx context.Context, req events.ALBTargetGroupRequest) (res events.ALBTargetGroupResponse, err error) {

		authHeader := req.Headers["authorization"]

		userID, err := auth.ParseAuthHeader(authHeader)
		if err != nil {
			res.StatusCode = http.StatusInternalServerError
			e := errors.Wrap(err, "cannot parse the access token")
			res.Body = e.Error()
			return
		}

		ctx = context.WithValue(ctx, "UserID", userID)

		res, err = h(ctx, req)
		return
	}
}