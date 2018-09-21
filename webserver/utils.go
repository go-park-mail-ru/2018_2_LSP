package webserver

import (
	"encoding/json"
	"net/http"
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

func writeJSONToStream(w http.ResponseWriter, p jsonConvertable) error {
	return json.NewEncoder(w).Encode(p)
}
