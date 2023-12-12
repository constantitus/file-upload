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
var Limited cache.Cache[string, bool]

func init() {
    gen = uuid.NewGen()

    UUID = cache.NewCache[string, UserData]()
    Limited = cache.NewCache[string, bool]()
    go func() {
        for {
            time.Sleep(time.Minute * 1)
            UUID.DeleteExpired()
            Limited.DeleteExpired()
        }
    }()
}

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
    http.SetCookie(*w, &http.Cookie{Name: "uuid", Path: "/", Expires: time.Now()})
    return false
}

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

func setUser(form *LoginForm) (cookie *http.Cookie) {
    id := gen.NewV4().String()
    var expires time.Time
    if form.Remember {
        UUID.Set(id, UserData{form.Username, form.admin},
            time.Duration(Conf.UUID_long_ttl) * time.Hour)
        expires = time.Now().Add(time.Duration(Conf.UUID_long_ttl) * time.Hour)
    } else {
        UUID.Set(id, UserData{form.Username, form.admin},
            time.Duration(Conf.UUID_def_ttl) * time.Hour)
        expires = time.Now().Add(time.Duration(Conf.UUID_def_ttl) * time.Hour)
    }
    

    cookie = &http.Cookie{
        Name: "uuid",
        Value: id,
        Expires: expires,
        Path: "/",
    }

    return
}
