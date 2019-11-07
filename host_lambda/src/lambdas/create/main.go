package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"log"
	"thumbnailr/app/create"
	"thumbnailr/repos_dynamodb"
	"thumbnailr/stores_s3"
)

var appHandler *create.Handler

func init() {
	sess, err := session.NewSession()
	if err != nil {
		log.Printf("cannot create new sdk session: %s", err)
		return
	}

	photoStore := stores_s3.NewPhotoStore(sess, "thumbnailr-photostore")
	thumbnailStore := stores_s3.NewThumbnailStore(sess, "thumbnailr-thumbnailstore")
	thumbnailRepo := repos_dynamodb.NewThumbnailRepo(sess)

	appHandler = create.NewHandler(photoStore, thumbnailStore, thumbnailRepo)
}

func lambdaHandler(evt events.SNSEvent) error {
	for _, r := range evt.Records {
		if err := handleRecord(&r); err != nil {
			return err
		}
	}
	return nil
}

func handleRecord(rec *events.SNSEventRecord) error {

	var in create.Input;
	msgBytes := []byte(rec.SNS.Message)
	if err := json.Unmarshal(msgBytes, &in); err != nil {
		return errors.Wrap(err,"cannot unmarshal the SNS message as create input")
	}

	out := appHandler.Handle(in)

	if out.IsUnexpected {
		return errors.Wrapf(out.Error, "unexpected error: %s", out.Message)
	}

	return nil
}

func main() {
	lambda.Start(lambdaHandler)
}
