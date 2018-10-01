package middlewares

import (
	"encoding/json"
	"net/http"
)

type apiError struct {
	Code    int
	Message string
}

func responseJSON(statusCode int, w http.ResponseWriter, p interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(p)
}
