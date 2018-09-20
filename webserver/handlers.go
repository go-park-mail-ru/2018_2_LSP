package webserver

import (
	"2018_2_LSP/user"
	"encoding/json"
	"net/http"
)

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
		// TODO добавить возврат ошибок по полям
		w.WriteHeader(http.StatusConflict)
		writeJSONToStream(w, apiError{2, err.Error()})
		return
	}

	err = writeJSONToStream(w, apiAuth{0, u.Token})
}
