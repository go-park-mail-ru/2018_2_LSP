package leaderboard

import (
	"github.com/go-park-mail-ru/2018_2_LSP/utils"
)

// Item Structure that stores cell for leaderboard
type Item struct {
	Username string
	Rating   int
}

// GetPage Returns 10 leaders for specified page
func GetPage(page int) ([]Item, error) {
	items := make([]Item, 0)

	rows, err := utils.Query("SELECT username, rating FROM users ORDER BY rating DESC LIMIT 10 OFFSET $1", page*10)
	if err != nil {
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var i Item
		err = rows.Scan(&i.Username, &i.Rating)
		if err != nil {
			return items, err
		}
		items = append(items, i)
	}

	return items, nil
}
