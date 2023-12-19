package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

    "github.com/go-pkgz/expirable-cache/v2"
	"nullprogram.com/x/uuid"
)

type UserData struct {
    Name string
    Rank bool
}

var gen *uuid.Gen

var UUID cache.Cache[string, UserData]

func init() {
    gen = uuid.NewGen()

    UUID = cache.NewCache[string, UserData]()
    go func() {
        for {
            time.Sleep(time.Minute * 1)
            UUID.DeleteExpired()
        }
    }()
}

// Checks the UUID from the cookie against the memory cache
func FromCookie(r *http.Request) (user UserData) {
    uuidCookie, err := r.Cookie("uuid")
    if err != nil {
        return
    }
    uuidString, _ := strings.CutPrefix(uuidCookie.String(), "uuid=")
    if uuidString == "" {
        return
    }

    value, got := UUID.Get(uuidString)
    if got {
        return value
    }
    return
}

func ClearFromCache(r *http.Request) {
    uuidCookie, err := r.Cookie("uuid")
    if err != nil {
        return
    }
    uuidString, _ := strings.CutPrefix(uuidCookie.String(), "uuid=")
    if uuidString == "" {
        return
    }

    UUID.Set(uuidString, UserData{}, time.Duration(1))
    return
}

// Checks the credentials against the database. If valid, it generates a UUID
// which it saves in the memory cache and in a cookie.
func CheckCredentials(form *LoginData, w *http.ResponseWriter) (valid bool) {
    var hash string
    hash, form.admin = QueryDB(form.Username)
    tmp := sha256.New()
    tmp.Write([]byte(form.Password))
    if hash == hex.EncodeToString(tmp.Sum(nil)) {
        http.SetCookie(*w, setUser(form))
        return true
    }
    // clear the cookie
    http.SetCookie(*w,
        &http.Cookie{Name: "uuid", Path: "/", Expires: time.Now()})
    return false
}

func setUser(form *LoginData) (cookie *http.Cookie) {
    id := gen.NewV4().String()
    var expires time.Time
    var ttl time.Duration
    if form.Remember {
        ttl = Config.UuidLongTTL
    } else {
        ttl = Config.UuidDefTTL
    }
    UUID.Set(id, UserData{form.Username, form.admin}, ttl)
    expires = time.Now().Add(ttl)
    

    cookie = &http.Cookie{
        Name: "uuid",
        Value: id,
        Expires: expires,
        Path: "/",
    }

    return
}

