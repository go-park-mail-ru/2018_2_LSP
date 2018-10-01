package webserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2018_2_LSP/user"
)

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
		responseJSON(http.StatusBadRequest, w, apiError{1, err.Error()})
		return
	}

	var u user.User
	errs := u.Auth(c)
	if errs != nil {
		responseJSON(http.StatusBadRequest, w, apiError{2, err.Error()})
		return
	}

	setAuthCookies(w, u.Token)
	responseJSON(http.StatusOK, w, apiAuth{0, u.Token})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var u user.User
	err := decoder.Decode(&u)
	if err != nil {
		responseJSON(http.StatusBadRequest, w, apiError{1, err.Error()})
		return
	}

	if errs := u.Register(); errs != nil {
		fields := make([]fieldError, 0)
		for _, e := range errs {
			if e, ok := e.(*user.RegisterError); ok {
				fields = append(fields, fieldError{e.Field, e.Message})
			}
		}
		responseJSON(http.StatusConflict, w, registerError{1, fields})
		return
	}

	setAuthCookies(w, u.Token)
	responseJSON(http.StatusOK, w, apiAuth{0, u.Token})
}
