package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

    "main/config"
    "main/cache"
)

var database *sql.DB

// Initialize the database
func Initialize() error {
    if err := Read(); err != nil {
        return err
    }

    // users table
    statement, err := database.Prepare( `CREATE TABLE IF NOT EXISTS users` +
        `(username TEXT PRIMARY KEY, password TEXT, rank INTEGER DEFAULT 0)`)
    if err != nil {
        return err
    }
    statement.Exec()

    // cache table
    statement, err = database.Prepare( `CREATE TABLE IF NOT EXISTS cache` +
        `(uuid TEXT PRIMARY KEY, user TEXT, rank INTEGER, expire INTEGER)`)
    if err != nil {
        return err
    }
    statement.Exec()

    parseCache()
    return nil
}

func Read() error {
    var err error
    database, err = sql.Open("sqlite3", config.DatabasePath)
    if err != nil {
        return err
    }
    return nil
}

// Query the database for username
func Query(username string) (hash string, admin bool) {
    rows, _ := database.Query(
        `SELECT password, rank FROM users WHERE username = ?`, username)
    defer rows.Close()

    if rows.Next() {
        rows.Scan(&hash, &admin)
    }
    return
}

func parseCache() {
    // Clear old entries
    _, err := database.Exec("DELETE FROM cache WHERE expire < ?",
        time.Now().UnixNano())
    if err != nil {
        log.Println(err)
    }

    // Read the cache table
    rows, _ := database.Query("SELECT * FROM cache")
    defer rows.Close()

    for rows.Next() {
        var data cache.Data
        var uuid string
        var tmp int64
        rows.Scan(&uuid, &data.Name, &data.Rank, &tmp)
        exp := time.Unix(0, tmp)
        if exp.After(exp) {
            continue
        }
        cache.Set(uuid, data, time.Until(exp))
    }
}

func StoreCache() {
    // Clear the table
    _, err := database.Exec("DELETE FROM cache")
    if err != nil {
        log.Println(err)
    }

    cache.DeleteExpired()
    for _, uuid := range cache.Keys() {
        data, got := cache.Get(uuid)
        if got {
            exp, _ := cache.GetExp(uuid)
            addCache(uuid, data, exp.UnixNano())
        }
    }
}

func addCache(uuid string, data cache.Data, exp int64) {
    _, err := database.Exec("INSERT INTO cache VALUES (?,?,?,?)",
        uuid, data.Name, data.Rank, int(exp))
    if err != nil {
        log.Println(err)
    }
}

func AddUser(username string, hash string, rank bool) error {
    _, err := database.Exec("INSERT INTO users VALUES (?,?,?)",
        username, hash, rank)
    if err != nil {
        return err
    }

    return nil
}

func RemoveUser(username string) error {
    _, err := database.Exec("DELETE FROM users WHERE username = (?)", username)
    if err != nil {
        return err
    }

    return nil
}


func QueryAll() (result []struct{user string; rank bool}) {
    rows, _ := database.Query(
        `SELECT username, rank FROM users`)
    defer rows.Close()

    for i := 0; rows.Next(); i++ {
        result = append(result, struct{user string; rank bool}{})
        rows.Scan(&result[i].user, &result[i].rank)
    }
    return
}
