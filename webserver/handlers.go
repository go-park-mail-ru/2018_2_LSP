package webserver

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2018_2_LSP/user"
	"github.com/thedevsaddam/govalidator"
)

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	signature, err := r.Cookie("signature")
	if err == nil {
		signature.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, signature)
	}

	headerPayload, err := r.Cookie("header.payload")
	if err == nil {
		headerPayload.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, headerPayload)
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	var c user.Credentials
	rules := govalidator.MapData{
		"email":    []string{"required", "between:4,25", "email"},
		"password": []string{"required", "alpha_space"},
	}

	opts := govalidator.Options{
		Request: r,
		Data:    &c,
		Rules:   rules,
	}
	v := govalidator.New(opts)
	if e := v.ValidateJSON(); len(e) > 0 {
		err := map[string]interface{}{"validationError": e}
		responseJSON(http.StatusBadRequest, w, err)
		return
	}

	var u user.User
	if err := u.Auth(c); err != nil {
		responseJSON(http.StatusBadRequest, w, apiError{2, err.Error()})
		return
	}

	setAuthCookies(w, u.Token)
	responseJSON(http.StatusOK, w, apiAuth{0, u.Token})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var u user.User
	rules := govalidator.MapData{
		"username":  []string{"required", "between:4,25"},
		"email":     []string{"required", "between:4,25", "email"},
		"password":  []string{"required", "alpha_space"},
		"firstname": []string{"alpha_space", "between:4,25"},
		"lastname":  []string{"alpha_space", "between:4,25"},
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

	if err := u.Register(); err != nil {
		responseJSON(http.StatusConflict, w, apiError{1, err.Error()})
		return
	}

	setAuthCookies(w, u.Token)
	responseJSON(http.StatusOK, w, apiAuth{0, u.Token})
}
