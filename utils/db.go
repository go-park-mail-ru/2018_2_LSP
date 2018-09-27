package utils

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB
var connectionOpened = false

// Query query sql request with args
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	if !connectionOpened {
		connStr := "host= " + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASS") + " dbname=" + os.Getenv("DB_DB") + " sslmode=disable"
		db, _ = sql.Open("postgres", connStr)
		connectionOpened = true
	}
	return db.Query(query, args...)
}
