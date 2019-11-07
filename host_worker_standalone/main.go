package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"thumbnailr/app"
	"thumbnailr/app/create"
	"thumbnailr/repos_dynamodb"
	"thumbnailr/stores_s3"
)

func subscribe(endPoint string, topicARN string) error {
	input := &sns.SubscribeInput{
		Endpoint: aws.String(endPoint),
		Protocol: aws.String("http"),
		TopicArn: aws.String(topicARN),
	}

	sess, err := session.NewSession()
	if err != nil {
		return errors.Wrap(err,"cannot create new aws sdk session")
	}

	svc := sns.New(sess)

	if _, err := svc.Subscribe(input); err != nil {
		return errors.Wrap(err,"cannot subscribe")
	}

	return nil
}

func confirmSubscription(subscribeURL string) error {
	if _, err := http.Get(subscribeURL); err != nil {
		return errors.Wrapf(err,"cannot http-get %s", subscribeURL)
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
			return errors.Wrap(err,"cannot unmarshal the create input msg")
		}

		out := appHandler.Handle(in)

		return outputToError("create",  out)
	}
}

func outputToError(msgHandlerName string, out app.Output) error {
	if out.IsUnexpected && out.Error != nil {
		return errors.Wrapf(out.Error, "message handler '%s' returned an unexpected error", msgHandlerName)
	}
	return nil
}


func main() {

	sess, err := session.NewSession()
	if err != nil {
		log.Printf("cannot create new sdk session: %s", err)
		return
	}

	photoStore := stores_s3.NewPhotoStore(sess, "thumbnailr-photostore")
	thumbnailStore := stores_s3.NewThumbnailStore(sess, "thumbnailr-thumbnailstore")
	thumbnailRepo := repos_dynamodb.NewThumbnailRepo(sess)

	appHandler := create.NewHandler(photoStore, thumbnailStore, thumbnailRepo)

	msgHandlers := map[string]MessageHandlerFunc {
		"create": createMessageHandler(appHandler),
	}

	http.HandleFunc("/", snsHandler(msgHandlers))
	log.Fatal(http.ListenAndServe(":9098", nil))
}
