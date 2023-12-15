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

// Checks the credentials against the database. If valid, it generates a UUID
// which it saves in the memory cache and in a cookie.
func CheckCredentials(form *LoginForm, w *http.ResponseWriter) (valid bool) {
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

func setUser(form *LoginForm) (cookie *http.Cookie) {
    id := gen.NewV4().String()
    var expires time.Time
    if form.Remember {
        UUID.Set(id, UserData{form.Username, form.admin},
            Config.UuidLongTTL * time.Hour)
        expires = time.Now().Add(Config.UuidLongTTL * time.Hour)
    } else {
        UUID.Set(id, UserData{form.Username, form.admin},
            Config.UuidDefTTL * time.Hour)
        expires = time.Now().Add(Config.UuidDefTTL * time.Hour)
    }
    

    cookie = &http.Cookie{
        Name: "uuid",
        Value: id,
        Expires: expires,
        Path: "/",
    }

    return
}
