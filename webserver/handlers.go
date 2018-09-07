package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shakal/user"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var u user.Credentials
	err := decoder.Decode(&u)
	p := apiAnswer{}
	if err != nil {
		p = apiAnswer{Code: 1, Message: "Error during parsing data"}
	} else {
		err = user.Auth(u)
		if err != nil {
			p = apiAnswer{Code: 1, Message: err.Error()}
		} else {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"foo": "bar",
				"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
			})
			tokenString, _ := token.SignedString([]byte("Secret"))
			p = apiAnswer{Code: 0, Message: tokenString}
		}
	}
	err = writeJSONToStream(w, p)
}
