package webserver

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2018_2_LSP/leaderboard"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pageStr, err := extractKey(r, "page")
	if err != nil {
		responseJSON(http.StatusBadRequest, w, apiError{2, err.Error()})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		responseJSON(http.StatusBadRequest, w, apiError{3, err.Error()})
		return
	}

	cells, err := leaderboard.GetPage(page)
	if err != nil {
		responseJSON(http.StatusBadRequest, w, apiError{4, err.Error()})
		return
	}

	responseJSON(http.StatusOK, w, cells)
}
