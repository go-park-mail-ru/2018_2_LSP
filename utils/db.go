package utils

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB
var connectionOpened = false

// Query query sql request with args
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	if !connectionOpened {
		connStr := "user=postgres password=root dbname=mytestdb sslmode=disable"
		db, _ = sql.Open("postgres", connStr)
		connectionOpened = true
	}
	return db.Query(query, args...)
}
