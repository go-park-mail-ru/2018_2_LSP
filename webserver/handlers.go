package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2018_2_LSP/user"
)

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	signature, err := r.Cookie("signature")
	if err == nil {
		signature.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, signature)
	}
	headerPayload, err := r.Cookie("header.payload")
	if err == nil {
		headerPayload.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, headerPayload)
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var c user.Credentials
	err := decoder.Decode(&c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSONToStream(w, apiError{1, err.Error()})
		return
	}

	u, err := user.Auth(c)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		writeJSONToStream(w, apiError{2, err.Error()})
		return
	}

	firstDot := strings.Index(u.Token, ".") + 1
	splitter := strings.Index(u.Token[firstDot:], ".") + firstDot
	fmt.Println(firstDot, splitter)
	cookieHeaderPayload := http.Cookie{
		Name:    "header.payload",
		Value:   u.Token[:splitter],
		Expires: time.Now().Add(30 * time.Minute),
	}
	cookieSignature := http.Cookie{
		Name:    "signature",
		Value:   u.Token[splitter+1:],
		Expires: time.Now().Add(720 * time.Hour),
	}
	http.SetCookie(w, &cookieHeaderPayload)
	http.SetCookie(w, &cookieSignature)

	err = writeJSONToStream(w, apiAuth{0, u.Token})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var u user.User
	err := decoder.Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSONToStream(w, apiError{1, err.Error()})
		return
	}

	u, err = user.Register(u)
	if err != nil {
		if err, ok := err.(*user.RegisterError); ok {
			fields := make([]fieldError, 0)
			switch err.Code() {
			case 1:
				fields = append(fields, fieldError{"username", err.Error()})
				fields = append(fields, fieldError{"email", err.Error()})
			case 2:
				fields = append(fields, fieldError{"username", err.Error()})
			case 3:
				fields = append(fields, fieldError{"username", err.Error()})
			}
			writeJSONToStream(w, registerError{2, fields})
			return
		}
		w.WriteHeader(http.StatusConflict)
		writeJSONToStream(w, apiError{2, err.Error()})
		return
	}

	err = writeJSONToStream(w, apiAuth{0, u.Token})
}
