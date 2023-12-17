package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB

// Initialize the database
func InitDB() {
    var err error
    Database, err = sql.Open("sqlite3", Config.DatabasePath)
    if err != nil {
        log.Panic(err)
    }

    // users table
    statement, err := Database.Prepare(
        "CREATE TABLE IF NOT EXISTS users (username TEXT PRIMARY KEY, password TEXT, rank INTEGER DEFAULT 0)")
    if err != nil {
        log.Panic(err)
    }
    statement.Exec()

    // cache table
    statement, err = Database.Prepare(
        "CREATE TABLE IF NOT EXISTS cache (uuid TEXT PRIMARY KEY, user TEXT, rank INTEGER, expire INTEGER)")
    if err != nil {
        log.Panic(err)
    }
    statement.Exec()
}

// Query the database for username
func QueryDB(username string) (hash string, admin bool) {
    rows, _ := Database.Query(
        "SELECT password, rank FROM users WHERE username = ?", username)
    defer rows.Close()

    if rows.Next() {
        rows.Scan(&hash, &admin)
    }
    return
}


func StoreCacheDB() {
    UUID.DeleteExpired()
    for _, uuid := range UUID.Keys() {
        data, got := UUID.Get(uuid)
        if got {
            exp, _ := UUID.GetExp(uuid)
            addCacheDB(uuid, data, exp.UnixNano())
        }
    }
}

func addCacheDB(uuid string, data UserData, exp int64) {
    _, err := Database.Exec("INSERT INTO cache VALUES (?,?,?,?)", uuid, data.Name, data.Rank, int(exp))
    if err != nil {
        log.Println(err)
    }
}


func ParseCacheDB() {
    // Clear old entries
    _, err := Database.Exec("DELETE FROM cache WHERE expire < ?", time.Now().UnixNano())
    if err != nil {
        log.Println(err)
    }

    // Read the cache table
    rows, _ := Database.Query("SELECT * FROM cache")
    defer rows.Close()

    for rows.Next() {
        var data UserData
        var uuid string
        var tmp int64
        rows.Scan(&uuid, &data.Name, &data.Rank, &tmp)
        exp := time.Unix(0, tmp)
        if exp.After(exp) {
            continue
        }
        UUID.Set(uuid, data, time.Until(exp))
    }

    // Clear the table
    _, err = Database.Exec("DELETE FROM cache")
    if err != nil {
        log.Println(err)
    }
}
