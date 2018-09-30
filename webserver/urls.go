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
	lt := mw.Chain(mw.Logging, mw.Tracing)
	return []urlHandler{
		urlHandler{"/", lt(mainHandler)},
		urlHandler{"/avatars", lt(avatarsHandler)},
	}
}
