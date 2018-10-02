package webserver

import (
	"net/http"

	mw "github.com/go-park-mail-ru/2018_2_LSP/webserver/middlewares"
)

type urlHandler struct {
	url     string
	handler http.HandlerFunc
}

func getUrlsAndHandlers() []urlHandler {
	lt := mw.Chain(mw.Cors, mw.Logging, mw.Tracing)
	return []urlHandler{
		urlHandler{"/auth", lt(mw.Post(authHandler))},
		urlHandler{"/register", lt(mw.Post(registerHandler))},
		urlHandler{"/logout", lt(mw.Delete(mw.Auth(logoutHandler)))},
	}
}
