package shared

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"net/http"
	"thumbnailr/app"
)

func OutputToResp(output *app.Output, res *events.APIGatewayProxyResponse) {
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