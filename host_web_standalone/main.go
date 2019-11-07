package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"thumbnailr/app"
	"thumbnailr/app/check_creation"
	"thumbnailr/app/get"
	"thumbnailr/app/request_creation"
	"thumbnailr/bus_sns"
	"thumbnailr/repos_dynamodb"
	"thumbnailr/stores_s3"
)

var auth app.Auth

func addAuth(auth *app.Auth, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		prefix := "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		token := strings.ReplaceAll(authHeader, prefix, "")

		userID, err := auth.ParseAuthHeader(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			e := errors.Wrap(err, "cannot parse the access token")
			w.Write([]byte(e.Error()))
			return
		}

		ctx := context.WithValue(r.Context(), "UserID", userID)
		r = r.WithContext(ctx)

		h(w,r)
	}
}

func wrapMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return addAuth(&auth,f)
}


func outputToResp(output *app.Output, w http.ResponseWriter) {
	if output.Success {
		w.WriteHeader(http.StatusOK)

		if data, err := json.Marshal(output.Data); err == nil {
			w.Write(data)

		} else {
			w.WriteHeader(http.StatusInternalServerError)
			e := errors.Wrap(err, "cannot marshal output data")
			w.Write([]byte(e.Error()))
		}

	} else {
		if output.IsUnexpected {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(output.Error.Error()))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(output.Message))
		}
	}
}

func checkCreationHandler(appHandler check_creation.Handler) http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		id := q.Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing the 'id' parameter"))
			return
		}

		in := check_creation.Input{
			UserID: r.Context().Value("UserID").(string),
			ThumbnailID: id,
		}

		out := appHandler.Handle(in)

		outputToResp(&out, w)
	}
	return wrapMiddleware(f)
}

func getHandler(appHandler get.Handler) http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		id := q.Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing the 'id' parameter"))
			return
		}

		in := get.Input{
			UserID: r.Context().Value("UserID").(string),
			ThumbnailID: id,
		}

		out := appHandler.Handle(in)

		outputToResp(&out, w)
	}
	return wrapMiddleware(f)
}

func requestCreationHandler(appHandler request_creation.Handler) http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		width, err := strconv.Atoi(q.Get("width"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing or the 'width' parameter is not an integer"))
			return
		}
		length, err := strconv.Atoi(q.Get("length"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing or the 'length' parameter is not an integer"))
			return
		}

		in := request_creation.Input{
			UserID: r.Context().Value("UserID").(string),
			PhotoID: q.Get("photoID"),
			Format: q.Get("format"),
			Width: width,
			Length: length,
		}

		out := appHandler.Handle(in)

		outputToResp(&out, w)
	}
	return wrapMiddleware(f)
}

func main() {

	sess, err := session.NewSession()
	if err != nil {
		return
	}

	auth = app.Auth{PrivateKey:"no-key"}

	cmdIssuer := bus_sns.NewCommandIssuer(sess)

	quotaRepo := repos_dynamodb.NewQuotaRepo(sess)

	thumbnailRepo := repos_dynamodb.NewThumbnailRepo(sess)

	thumbnailStore := stores_s3.NewThumbnailStore(sess, "thumbnailr-thumbnailstore")

	http.HandleFunc("/check_creation", checkCreationHandler(check_creation.Handler{
		ThumbnailRepo: thumbnailRepo,
		}))

	http.HandleFunc("/get", getHandler(get.Handler{
		ThumbnailRepo: thumbnailRepo,
		ThumbnailStore: thumbnailStore}))

	http.HandleFunc("/request_creation", requestCreationHandler(request_creation.Handler{
		ThumbnailRepo: thumbnailRepo,
		QuotaRepo: quotaRepo,
		CommandIssuer: cmdIssuer}))

	log.Fatal(http.ListenAndServe(":9097", nil))
}
