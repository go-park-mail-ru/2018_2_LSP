package webserver

import (
	mw "2018_2_LSP/webserver/middlewares"
	"net/http"
)

type urlHandler struct {
	url     string
	handler http.HandlerFunc
}

func getUrlsAndHandlers() []urlHandler {
	lt := mw.Chain(mw.Logging, mw.Tracing)
	return []urlHandler{
		urlHandler{"/", lt(mainHandler)},
		urlHandler{"/auth", lt(authHandler)},
	}
}
