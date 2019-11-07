package bus_sns

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/pkg/errors"
)

type SNSCommandIssuer struct {
	SNSClient *sns.SNS
}

func NewCommandIssuer(sess *session.Session) *SNSCommandIssuer {
	return &SNSCommandIssuer{
		SNSClient: sns.New(sess),
	}
}

func(ci *SNSCommandIssuer) Send(cmd interface{}) error {

	msg, err := json.Marshal(cmd)
	if err != nil {
		return errors.Wrap(err,"cannot marshal command")
	}

	input := &sns.PublishInput{
		Message:  aws.String(string(msg)),
		TopicArn: aws.String("thumbnailr-creation-requests"),
	}

	if _, err := ci.SNSClient.Publish(input); err != nil {
		return errors.Wrap(err,"cannot publish message")
	}

	return nil
}