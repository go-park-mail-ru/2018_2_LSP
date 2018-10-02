package webserver

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2018_2_LSP/leaderboard"
	"github.com/thedevsaddam/govalidator"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	rules := govalidator.MapData{
		"page": []string{"required", "numeric"},
	}

	opts := govalidator.Options{
		Request:         r,
		Rules:           rules,
		RequiredDefault: true,
	}
	v := govalidator.New(opts)
	if e := v.Validate(); len(e) > 0 {
		err := map[string]interface{}{"validationError": e}
		responseJSON(http.StatusBadRequest, w, err)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query()["page"][0])

	cells, err := leaderboard.GetPage(page)
	if err != nil {
		responseJSON(http.StatusBadRequest, w, apiError{4, err.Error()})
		return
	}

	responseJSON(http.StatusOK, w, cells)
}
