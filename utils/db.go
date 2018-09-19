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
		connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
		db, _ = sql.Open("postgres", connStr)
		connectionOpened = true
	}
	return db.Query(query, args)
}
