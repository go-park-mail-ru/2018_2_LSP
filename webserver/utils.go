package webserver

import (
	"encoding/json"
	"errors"
	"net/http"
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
