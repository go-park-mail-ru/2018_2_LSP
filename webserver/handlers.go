package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shakal/user"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var c user.Credentials
	err := decoder.Decode(&c)
	p := apiAnswer{}
	if err != nil {
		p = apiAnswer{Code: 1, Message: "Error during parsing data"}
	} else {
		u, err := user.Auth(c)
		if err != nil {
			p = apiAnswer{Code: 1, Message: err.Error()}
		} else {
			p = apiAnswer{Code: 0, Message: u.token}
		}
	}
	err = writeJSONToStream(w, p)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var u user.User
	err := decoder.Decode(&u)
	p := apiAnswer{}
	if err != nil {
		p = apiAnswer{Code: 1, Message: "Error during parsing data"}
	} else {
		u, err := user.Register(u)
		if err != nil {
			p = apiAnswer{Code: 1, Message: err.Error()}
		} else {
			p = apiAnswer{Code: 0, Message: u.token}
		}
	}
	err = writeJSONToStream(w, p)
}
