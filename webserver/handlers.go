package webserver

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-park-mail-ru/2018_2_LSP/user"
	"github.com/go-park-mail-ru/2018_2_LSP/utils"
)

func handlePutRequest(w http.ResponseWriter, r *http.Request, claims jwt.MapClaims) {
	var data map[string]string

	firstname, err := extractKey(r, "firstname")
	if err == nil {
		data["firstname"] = firstname
	}

	lastname, err := extractKey(r, "lastname")
	if err == nil {
		data["lastname"] = lastname
	}

	newPassword, err := extractKey(r, "newPassword")
	if err == nil {
		oldPassword, err := extractKey(r, "oldPassword")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJSONToStream(w, apiError{3, "Please, specify old password"})
			return
		}
		row, err := utils.Query("SELECT password FROM users WHERE id = $1", claims["id"].(int))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJSONToStream(w, apiError{3, err.Error()})
			return
		}
		var hashedPassword string
		row.Next()
		err = row.Scan(&hashedPassword)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJSONToStream(w, apiError{3, err.Error()})
			return
		}
		if !utils.ComparePasswords(hashedPassword, oldPassword) {
			w.WriteHeader(http.StatusBadRequest)
			writeJSONToStream(w, apiError{3, "Wrong old password"})
			return
		}
		data["password"] = newPassword
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
	request += "WHERE id = $1 RETURNING firstname, lastname, email, username"
	rows, err := utils.Query(request, claims["id"].(int))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSONToStream(w, apiError{3, err.Error()})
		return
	}
	var u user.User
	rows.Next()
	err = rows.Scan(&u.FirstName, &u.LastName, &u.Email, &u.Username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSONToStream(w, apiError{3, err.Error()})
		return
	}
	writeJSONToStream(w, u)
}

func handleGetRequest(w http.ResponseWriter, r *http.Request, claims jwt.MapClaims) {
	u, err := user.GetUserByID(claims["id"].(int))
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
