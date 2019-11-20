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

	arn *string
}

func NewCommandIssuer(sess *session.Session) (*SNSCommandIssuer, error) {

	client := sns.New(sess)

	// this is done only for getting the topic ARN from the topic name
	out, err := client.CreateTopic(&sns.CreateTopicInput{
		Name: aws.String("thumbnailr-creation-requests"),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create topic")
	}

	return &SNSCommandIssuer{
		SNSClient: client,

		arn: out.TopicArn,
	}, nil
}

func (ci *SNSCommandIssuer) Send(cmd interface{}) error {

	msg, err := json.Marshal(cmd)
	if err != nil {
		return errors.Wrap(err, "cannot marshal command")
	}

	input := &sns.PublishInput{
		Message:  aws.String(string(msg)),
		TopicArn: ci.arn,
	}

	if _, err := ci.SNSClient.Publish(input); err != nil {
		return errors.Wrap(err, "cannot publish message")
	}

	return nil
}
