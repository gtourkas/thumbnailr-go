package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/handlerfunc"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"log"
	"net/http"
)

var handlerFuncLambda *handlerfunc.HandlerFuncAdapter

func init() {
	privateKey := []byte("no-key")

	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// generate jwt access token
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate(privateKey, jwt.SigningMethodHS512))

	clientStore := store.NewClientStore()
	clientStore.Set("web", &models.Client{
		ID:     "web",
		Secret: "no-secret",
	})
	manager.MapClientStorage(clientStore)

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetClientInfoHandler(func(r *http.Request) (clientID, clientSecret string, err error) {
		clientID = "web"
		clientSecret = "no-secret"
		return
	})
	srv.SetClientAuthorizedHandler(func(clientID string, grant oauth2.GrantType) (allowed bool, err error) {
		allowed = true
		return
	})
	srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		if username == "test" && password == "test" {
			userID = "testuser"
			log.Printf("successful pwd auth for username %s", username)
		} else {
			log.Printf("failed pwd auth for username %s", username)
		}
		return
	})

	handlerFuncLambda = handlerfunc.New(func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return handlerFuncLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
