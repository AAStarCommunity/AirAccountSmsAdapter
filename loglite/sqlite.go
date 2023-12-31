package loglite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/helpers/toto"
	"time"
)

const DB_FILE = "log.db"

func LogInfo(info string, msg *string) {
	log.Info(info, toto.V{"data": *msg})
	db, err := sql.Open("sqlite3", DB_FILE)
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	if err != nil {
		log.Fatal(err)
	}
	insertSQL := "INSERT INTO logs (msg, log_at) VALUES (?, ?)"
	_, err = db.Exec(insertSQL, info+*msg, time.Now())
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	db, err := sql.Open("sqlite3", DB_FILE)

	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	// Create a table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		msg TEXT,
		log_at DateTime
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

}
