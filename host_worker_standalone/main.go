package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"thumbnailr/app"
	"thumbnailr/app/create"
	"thumbnailr/repos_dynamodb"
	"thumbnailr/stores_s3"
)

func subscribe(sess *session.Session, endPoint string, topicARN string) error {
	input := &sns.SubscribeInput{
		Endpoint: aws.String(endPoint),
		Protocol: aws.String("http"),
		TopicArn: aws.String(topicARN),
	}

	svc := sns.New(sess)

	if _, err := svc.Subscribe(input); err != nil {
		return errors.Wrap(err, "cannot subscribe")
	}

	return nil
}

func confirmSubscription(subscribeURL string) error {
	if _, err := http.Get(subscribeURL); err != nil {
		return errors.Wrapf(err, "cannot http-get %s", subscribeURL)
	}
	return nil
}

func snsHandler(msgHandlers map[string]MessageHandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("cannot read request body: %s", err)
			return
		}

		var f interface{}
		err = json.Unmarshal(body, &f)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("cannot unmarshal request body: %s", err)
			return
		}

		data := f.(map[string]interface{})

		if data["Type"].(string) == "SubscriptionConfirmation" {

			subscribeURL := data["SubscribeURL"].(string)

			// confirm subscription
			if err := confirmSubscription(subscribeURL); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("cannot confirm subscription: %s", err)
				return
			}

			log.Print("subscription confirmed")
			w.WriteHeader(http.StatusOK)
			return

		} else if data["Type"].(string) == "Notification" {

			subj := data["Subject"].(string)
			msg := data["Message"].(string)

			if h, ok := msgHandlers[subj]; ok {
				if err := h(msg); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}

			w.WriteHeader(http.StatusOK)
			return

		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

type MessageHandlerFunc func(msg string) error

func createMessageHandler(appHandler *create.Handler) MessageHandlerFunc {
	return func(msg string) error {
		var in create.Input
		if err := json.Unmarshal([]byte(msg), &in); err != nil {
			return errors.Wrap(err, "cannot unmarshal the create input msg")
		}

		out := appHandler.Handle(in)

		return outputToError("create", out)
	}
}

func outputToError(msgHandlerName string, out app.Output) error {
	if out.IsUnexpected && out.Error != nil {
		return errors.Wrapf(out.Error, "message handler '%s' returned an unexpected error", msgHandlerName)
	}
	return nil
}

func main() {
	endPoint := os.Getenv("TN_WORKER_ENDPOINT")
	if endPoint == "" {
		endPoint = ":9098"
	}

	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("cannot create new sdk session: %s", err)
	}

	stss := sts.New(sess)
	cid, err := stss.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatalf("cannot get the caller identity: %s", err)
	}

	// this is done only for getting the topic ARN from the topic name
	snss := sns.New(sess)
	snso, err := snss.CreateTopic(&sns.CreateTopicInput{
		Name: aws.String("thumbnailr-creation-requests"),
	})
	if err != nil {
		log.Fatalf("cannot create the topic: %s", err)
	}
	snsTopicARN := *snso.TopicArn

	photoStore := stores_s3.NewPhotoStore(sess, fmt.Sprintf("thumbnailr-photostore-%s", *cid.Account))
	thumbnailStore := stores_s3.NewThumbnailStore(sess, fmt.Sprintf("thumbnailr-thumbnailstore-%s", *cid.Account))
	thumbnailRepo := repos_dynamodb.NewThumbnailRepo(sess)

	appHandler := create.NewHandler(photoStore, thumbnailStore, thumbnailRepo)

	msgHandlers := map[string]MessageHandlerFunc{
		"create": createMessageHandler(appHandler),
	}

	go func() {
		http.HandleFunc("/", snsHandler(msgHandlers))
	}()

	if err := subscribe(sess, endPoint, snsTopicARN); err != nil {
		log.Fatalf("cannot subscribe to SNS: %s", err)
	}

	log.Fatal(http.ListenAndServe(endPoint, nil))
}
