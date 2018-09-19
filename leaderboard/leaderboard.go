package leaderboard

import (
	"2018_2_LSP/utils"
)

// Item Structure that stores cell for leaderboard
type Item struct {
	Username string
	Rating   int
}

// GetPage Returns 10 leaders for specified page
func GetPage(page int) ([]Item, error) {
	items := make([]Item, 10)
	rows, err := utils.Query("SELECT username, rating FROM users ORDER BY rating LIMIT 10 OFFSET $1", page*10)
	if err != nil {
		return items, err
	}
	for i := range items {
		err = rows.Scan(&items[i])
		if err != nil {
			return items, nil
		}
	}

	return items, nil
}
