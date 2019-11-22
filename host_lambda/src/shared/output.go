package shared

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"net/http"
	"thumbnailr/app"
)

func OutputToApiGatewayResp(output *app.Output, res *events.APIGatewayProxyResponse) {
	res.Headers = map[string]string{"content-type":  "application/json"}
	if output.Success {
		res.StatusCode = http.StatusOK

		if data, err := json.Marshal(output.Data); err == nil {
			res.Body = string(data)

		} else {
			res.StatusCode = http.StatusInternalServerError
			err := errors.Wrap(err, "cannot marshal output data")
			res.Body = err.Error()
		}

	} else {
		if output.IsUnexpected {
			res.StatusCode = http.StatusInternalServerError
			res.Body = output.Error.Error()
		} else {
			res.StatusCode = http.StatusBadRequest
			res.Body = output.Message
		}
	}
}

func OutputToAlbTargetGroupResp(output *app.Output, res *events.ALBTargetGroupResponse) {
	res.Headers = map[string]string{"content-type":  "application/json"}
	res.IsBase64Encoded = false
	if output.Success {
		res.StatusCode = http.StatusOK

		if data, err := json.Marshal(output.Data); err == nil {
			res.Body = string(data)

		} else {
			res.StatusCode = http.StatusInternalServerError
			err := errors.Wrap(err, "cannot marshal output data")
			res.Body = err.Error()
		}

	} else {
		if output.IsUnexpected {
			res.StatusCode = http.StatusInternalServerError
			res.Body = output.Error.Error()
		} else {
			res.StatusCode = http.StatusBadRequest
			res.Body = output.Message
		}
	}

	res.StatusDescription = fmt.Sprintf("%d %s", res.StatusCode, http.StatusText(res.StatusCode))
}