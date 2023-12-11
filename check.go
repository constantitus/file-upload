package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"nullprogram.com/x/uuid"
)

const DefaultTTL = time.Hour * 2
const RememberTTL = time.Hour * 24

var UUID struct {
    Gen *uuid.Gen
    Id *ttlcache.Cache[string, string]
    Admin *ttlcache.Cache[string, bool]
}

func init() {
    UUID.Gen = uuid.NewGen()

    UUID.Id = ttlcache.New[string, string]()
    UUID.Admin = ttlcache.New[string, bool]()
    go UUID.Id.Start()
    go UUID.Admin.Start()
}

func setUser(form *LoginForm) *http.Cookie {
    id := UUID.Gen.NewV4().String()
    var expires time.Time
    if form.Remember {
        UUID.Id.Set(id, form.Username, RememberTTL)
        UUID.Admin.Set(form.Username, form.admin, RememberTTL)
        expires = time.Now().Add(RememberTTL)
    } else {
        UUID.Id.Set(id, form.Username, DefaultTTL)
        UUID.Admin.Set(form.Username, form.admin, DefaultTTL)
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

func CheckCredentials(form *LoginForm, w *http.ResponseWriter) bool {
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

    get := UUID.Id.Get(uuidString)
    if get != nil {
        fmt.Println(get)
        return get.Value()
    }
    return
}
