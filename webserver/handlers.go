package webserver

import (
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-park-mail-ru/2018_2_LSP/user"
	"github.com/go-park-mail-ru/2018_2_LSP/utils"
	"github.com/gorilla/context"
	"github.com/thedevsaddam/govalidator"
)

func avatarsHandler(w http.ResponseWriter, r *http.Request) {
	claims := context.Get(r, "claims").(jwt.MapClaims)

	rules := govalidator.MapData{
		"file:file": []string{"required", "ext:jpg,png", "size:300000", "mime:image/jpg,image/png"},
	}

	opts := govalidator.Options{
		Request: r,
		Rules:   rules,
	}
	v := govalidator.New(opts)
	if e := v.Validate(); len(e) > 0 {
		err := map[string]interface{}{"validationError": e}
		responseJSON(http.StatusBadRequest, w, err)
		return
	}

	file, handle, err := r.FormFile("file")
	if err != nil {
		responseJSON(http.StatusBadRequest, w, apiError{1, err.Error()})
		return
	}
	defer file.Close()

	var u user.User
	u.ID = int(claims["id"].(float64))

	err = saveFile(file, handle, u.ID)
	if err != nil {
		responseJSON(http.StatusBadRequest, w, apiError{1, err.Error()})
		return
	}
	response := avatarUpload{URL: "/avatars/" + strconv.Itoa(u.ID) + "_" + handle.Filename}
	responseJSON(http.StatusOK, w, response)
	return
}

func handlePutRequest(w http.ResponseWriter, r *http.Request, claims jwt.MapClaims) {
	var u user.User
	rules := govalidator.MapData{
		"firstname": []string{"between:4,25"},
		"lastname":  []string{"between:4,25"},
		"password":  []string{"alpha_space"},
	}

	opts := govalidator.Options{
		Request: r,
		Data:    &u,
		Rules:   rules,
	}
	v := govalidator.New(opts)
	if e := v.ValidateJSON(); len(e) > 0 {
		err := map[string]interface{}{"validationError": e}
		responseJSON(http.StatusBadRequest, w, err)
		return
	}

	data := make(map[string]string)
	if len(u.FirstName) > 0 {
		data["first_name"] = u.FirstName
	}
	if len(u.LastName) > 0 {
		data["last_name"] = u.LastName
	}

	u.ID = int(claims["id"].(float64))

	if len(u.Password) > 0 {
		if len(u.OldPassword) == 0 {
			responseJSON(http.StatusBadRequest, w, apiError{3, "Please, specify old password"})
			return
		}
		isValid, err := user.ValidateUserPassword(u.OldPassword, u.ID)
		if err != nil {
			responseJSON(http.StatusBadRequest, w, apiError{3, err.Error()})
			return
		}
		if !isValid {
			responseJSON(http.StatusBadRequest, w, apiError{3, "Wrong old password"})
			return
		}
		data["password"] = utils.HashAndSalt([]byte(u.Password))
	}

	if len(data) == 0 {
		responseJSON(http.StatusBadRequest, w, apiError{3, "Empty request"})
		return
	}

	request := "UPDATE users SET "

	for k, v := range data {
		request += k + "='" + v + "',"
	}
	request = request[:len(request)-1]
	request += " WHERE id = $1 RETURNING first_name, last_name, email, username"
	rows, err := utils.Query(request, u.ID)
	if err != nil {
		responseJSON(http.StatusBadRequest, w, apiError{3, err.Error()})
		return
	}
	rows.Next()
	err = rows.Scan(&u.FirstName, &u.LastName, &u.Email, &u.Username)
	if err != nil {
		responseJSON(http.StatusBadRequest, w, apiError{3, err.Error()})
		return
	}
	responseJSON(http.StatusOK, w, u)
}

func handleGetRequest(w http.ResponseWriter, r *http.Request, claims jwt.MapClaims) {
	id := int(claims["id"].(float64))
	u, err := user.GetUserByID(id)
	if err != nil {
		responseJSON(http.StatusBadRequest, w, apiError{1, err.Error()})
		return
	}
	responseJSON(http.StatusOK, w, u)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	claims := context.Get(r, "claims").(jwt.MapClaims)
	switch r.Method {
	case http.MethodPut:
		handlePutRequest(w, r, claims)
	case http.MethodGet:
		handleGetRequest(w, r, claims)
	}
}
