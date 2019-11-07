package app

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

type Auth struct {
	PrivateKey string
}

func (a *Auth) ParseAuthHeader(authHeader string) (userID string, err error)  {
	prefix := "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("the authorization header should start with 'Bearer '")
	}
	bearerToken := strings.ReplaceAll(authHeader, prefix, "")

	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.PrivateKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID = claims["sub"].(string)
	}
	return
}
