package webserver

import (
	"encoding/json"
	"net/http"
)

type jsonConvertable interface {
}

type apiAnswer struct {
	Code    int
	Message string
}

func writeJSONToStream(w http.ResponseWriter, p jsonConvertable) error {
	return json.NewEncoder(w).Encode(p)
}
