package webserver

import (
	"encoding/json"
	"errors"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

type jsonConvertable interface {
}

type apiError struct {
	Code    int
	Message string
}

type apiAuth struct {
	Code  int
	Token string
}

func extractKey(r *http.Request, key string) (string, error) {
	keys, ok := r.URL.Query()[key]

	if !ok || len(keys[0]) < 1 {
		return "", errors.New("Url Param " + key + " is missing")
	}
	return keys[0], nil
}

func writeJSONToStream(w http.ResponseWriter, p jsonConvertable) error {
	return json.NewEncoder(w).Encode(p)
}

func checkAuth(r *http.Request) (jwt.MapClaims, error) {
	signature, err := r.Cookie("signature")
	if err != nil {
		return nil, err
	}

	headerPayload, err := r.Cookie("header.payload")
	if err != nil {
		return nil, err
	}

	tokenString := headerPayload.Value + "." + signature.Value
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("HeAdfasdf3ref&^%$Dfrtgauyhia"), nil
	})

	return claims, err
}
