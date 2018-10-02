package webserver

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
)

type apiError struct {
	Code    int
	Message string
}

type apiAuth struct {
	Code  int
	Token string
}

type avatarUpload struct {
	URL string
}

func responseJSON(statusCode int, w http.ResponseWriter, p interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(p)
}

func saveFile(file multipart.File, handle *multipart.FileHeader, id int) error {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("/go/src/github.com/go-park-mail-ru/2018_2_LSP/avatars/"+strconv.Itoa(id)+"_"+handle.Filename, data, 0666)
	if err != nil {
		return err
	}

	return nil
}
