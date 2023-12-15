package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

// Initialize the database
func InitDB() {
    var err error
    database, err = sql.Open("sqlite3", Config.DatabasePath)
    if err != nil {
        log.Panic(err)
    }
    statement, err := database.Prepare(
        "CREATE TABLE IF NOT EXISTS users (username TEXT PRIMARY KEY, password TEXT, rank INTEGER DEFAULT 0)")
    if err != nil {
        log.Panic(err)
    }
    statement.Exec()
}

// Query the database for username
func QueryDB(username string) (hash string, admin bool) {
    rows, _ := database.Query(
        "SELECT password, rank FROM users WHERE username = ?", username)

    if rows.Next() {
        rows.Scan(&hash, &admin)
    }
    return
}
