package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

const DBfile = "./database.db"

var (
    database *sql.DB
)

func init() {
    var err error
    database, err = sql.Open("sqlite3", DBfile)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    statement, _ := database.Prepare(
        "CREATE TABLE IF NOT EXISTS users (username TEXT PRIMARY KEY, password TEXT, rank INTEGER DEFAULT 0)")
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
