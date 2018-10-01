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
	c := mw.Chain(mw.Cors, mw.Auth, mw.Logging, mw.Tracing)
	return []urlHandler{
		urlHandler{"/", c(mainHandler)},
	}
}
