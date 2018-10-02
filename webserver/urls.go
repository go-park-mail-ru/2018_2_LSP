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
		urlHandler{"/", lt(mw.GetPut(mw.Auth(mainHandler)))},
		urlHandler{"/avatar", lt(mw.Post(mw.Auth(avatarsHandler)))},
	}
}
