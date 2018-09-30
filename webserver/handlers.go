package webserver

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2018_2_LSP/leaderboard"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := checkAuth(r)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		writeJSONToStream(w, apiError{1, err.Error()})
		return
	}

	pageStr, err := extractKey(r, "page")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSONToStream(w, apiError{2, err.Error()})
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSONToStream(w, apiError{3, err.Error()})
		return
	}

	cells, err := leaderboard.GetPage(page)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSONToStream(w, apiError{4, err.Error()})
		return
	}
	err = writeJSONToStream(w, cells)
}
