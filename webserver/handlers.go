package webserver

import (
	"2018_2_LSP/leaderboard"
	"net/http"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cells, err := leaderboard.GetPage(0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSONToStream(w, apiError{1, err.Error()})
		return
	}
	err = writeJSONToStream(w, cells)
}
