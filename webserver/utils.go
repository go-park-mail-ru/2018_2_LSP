package webserver

import (
	"encoding/json"
	"errors"
	"net/http"
)

type apiError struct {
	Code    int
	Message string
}

func extractKey(r *http.Request, key string) (string, error) {
	keys, ok := r.URL.Query()[key]

	if !ok || len(keys[0]) < 1 {
		return "", errors.New("Url Param " + key + " is missing")
	}
	return keys[0], nil
}

func responseJSON(statusCode int, w http.ResponseWriter, p interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(p)
}
