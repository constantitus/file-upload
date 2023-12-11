package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func init() {
    var err error
    database, err = sql.Open("sqlite3", Conf.Database_path)
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


func QueryDB(username string) (hash string, admin bool) {
    rows, _ := database.Query(
        "SELECT password, rank FROM users WHERE username = ?", username)

    if rows.Next() {
        rows.Scan(&hash, &admin)
    }
    return
}
