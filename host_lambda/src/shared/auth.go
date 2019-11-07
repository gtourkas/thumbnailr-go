package shared

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"thumbnailr/app"
)


func AddAuth(auth *app.Auth, h ApiGatewayHandler) ApiGatewayHandler {
	return func(req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
		authHeader := req.Headers["Authorization"]
		prefix := "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			res.StatusCode = http.StatusForbidden
			return
		}
		token := strings.ReplaceAll(authHeader, prefix, "")

		userID, err := auth.ParseAuthHeader(token)
		if err != nil {
			res.StatusCode = http.StatusInternalServerError
			e := errors.Wrap(err, "cannot parse the access token")
			res.Body = e.Error()
			return
		}

		req.RequestContext.Identity.User = userID

		res, err = h(req)
		return
	}
}