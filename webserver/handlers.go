package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-park-mail-ru/2018_2_LSP/user"
	"github.com/go-park-mail-ru/2018_2_LSP/utils"
)

func handlePutRequest(w http.ResponseWriter, r *http.Request, claims jwt.MapClaims) {
	decoder := json.NewDecoder(r.Body)
	var u user.User
	err := decoder.Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSONToStream(w, apiError{1, err.Error()})
		return
	}

	data := make(map[string]string)

	if len(u.FirstName) > 0 {
		data["first_name"] = u.FirstName
	}

	if len(u.LastName) > 0 {
		data["last_name"] = u.LastName
	}

	if len(u.Password) > 0 {
		if len(u.OldPassword) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			writeJSONToStream(w, apiError{3, "Please, specify old password"})
			return
		}
		isValid, err := user.ValidateUserPassword(u.OldPassword, int(claims["id"].(float64)))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJSONToStream(w, apiError{3, err.Error()})
			return
		}
		if !isValid {
			w.WriteHeader(http.StatusBadRequest)
			writeJSONToStream(w, apiError{3, "Wrong old password"})
			return
		}
		data["password"] = u.Password
	}

	if len(data) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		writeJSONToStream(w, apiError{3, "Empty request"})
		return
	}

	request := "UPDATE users SET "

	for k, v := range data {
		request += k + "=" + v + ","
	}
	request = request[:len(request)-1]
	request += " WHERE id = $1 RETURNING first_name, last_name, email, username"
	fmt.Println(int(claims["id"].(float64)), request)
	rows, err := utils.Query(request, int(claims["id"].(float64)))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSONToStream(w, apiError{3, err.Error()})
		return
	}
	rows.Next()
	err = rows.Scan(&u.FirstName, &u.LastName, &u.Email, &u.Username)
	u.Password = ""
	u.OldPassword = ""
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSONToStream(w, apiError{3, err.Error()})
		return
	}
	writeJSONToStream(w, u)
}

func handleGetRequest(w http.ResponseWriter, r *http.Request, claims jwt.MapClaims) {
	id := int(claims["id"].(float64))
	u, err := user.GetUserByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSONToStream(w, apiError{3, err.Error()})
		return
	}
	writeJSONToStream(w, u)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPut && r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeJSONToStream(w, apiError{1, "Method not allowed"})
		return
	}

	claims, err := checkAuth(r)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		writeJSONToStream(w, apiError{2, err.Error()})
		return
	}

	switch r.Method {
	case http.MethodPut:
		handlePutRequest(w, r, claims)
	case http.MethodGet:
		handleGetRequest(w, r, claims)
	}
}
