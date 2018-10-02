package webserver

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type jsonConvertable interface {
}

type apiError struct {
	Code    int
	Message string
}

type registerError struct {
	Code   int
	Fields []fieldError
}

type fieldError struct {
	Field   string
	Message string
}

type apiAuth struct {
	Code  int
	Token string
}

func responseJSON(statusCode int, w http.ResponseWriter, p interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(p)
}

func setAuthCookies(w http.ResponseWriter, tokenString string) {
	firstDot := strings.Index(tokenString, ".") + 1
	secondDot := strings.Index(tokenString[firstDot:], ".") + firstDot
	cookieHeaderPayload := http.Cookie{
		Name:    "header.payload",
		Value:   tokenString[:secondDot],
		Expires: time.Now().Add(30 * time.Minute),
		Secure:  true,
	}
	cookieSignature := http.Cookie{
		Name:     "signature",
		Value:    tokenString[secondDot+1:],
		Expires:  time.Now().Add(720 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookieHeaderPayload)
	http.SetCookie(w, &cookieSignature)
}
