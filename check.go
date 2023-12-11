package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"nullprogram.com/x/uuid"
)

var DefaultTTL = time.Hour * time.Duration(Conf.Default_ttl)
var RememberTTL = time.Hour * time.Duration(Conf.Rememberme_ttl)

type UserData struct {
    User string
    Rank bool
}

var gen *uuid.Gen
var UUID *ttlcache.Cache[string, UserData]

func init() {
    gen = uuid.NewGen()

    UUID = ttlcache.New[string, UserData]()
    go UUID.Start()
}

func setUser(form *LoginForm) *http.Cookie {
    id := gen.NewV4().String()
    var expires time.Time
    if form.Remember {
        UUID.Set(id, UserData{form.Username, form.admin}, RememberTTL)
        expires = time.Now().Add(RememberTTL)
    } else {
        UUID.Set(id, UserData{form.Username, form.admin}, DefaultTTL)
        expires = time.Now().Add(DefaultTTL)
    }

    cookie := http.Cookie{
        Name: "uuid",
        Value: id,
        Expires: expires,
        Path: "/",
    }

    return &cookie
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
    http.SetCookie(*w, &http.Cookie{Name: "uuid", Path: "/"})
    return false
}

func FromCookie(r *http.Request) (user string) {
    uuidCookie, err := r.Cookie("uuid")
    if err != nil {
        return
    }
    uuidString, _ := strings.CutPrefix(uuidCookie.String(), "uuid=")
    if uuidString == "" {
        return
    }

    get := UUID.Get(uuidString)
    if get != nil {
        return get.Value().User
    }
    return
}
