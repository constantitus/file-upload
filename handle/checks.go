package handle

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"main/cache"
	"main/config"
	"main/db"
)

// Checks the credentials against the database. If valid, it generates a cache.UUID
// which it saves in the memory cache and in a cookie.
func checkCreds(form *loginData, w *http.ResponseWriter) (valid bool) {
    var hash string
    hash, form.admin = db.Query(form.Username)
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

func setUser(form *loginData) (cookie *http.Cookie) {
    id := gen.NewV4().String()
    var expires time.Time
    var ttl time.Duration
    if form.Remember {
        ttl = config.UuidLongTTL
    } else {
        ttl = config.UuidDefTTL
    }
    cache.UUID.Set(id, cache.Data{Name: form.Username, Rank: form.admin}, ttl)
    expires = time.Now().Add(ttl)
    
    cookie = &http.Cookie{
        Name: "uuid",
        Value: id,
        Expires: expires,
        Path: "/",
    }
    return
}

func getClientIP(r *http.Request) (ip string) {
    ip = r.Header.Get("X-Real-Ip")
    if ip == "" {
        ip = r.Header.Get("X-Forwarded-For")
    }
    if ip == "" {
        ip = r.RemoteAddr
    }
    return
}

// Checks the cache.UUID from the cookie against the memory cache
func fromCookie(r *http.Request) (user cache.Data) {
    uuidCookie, err := r.Cookie("uuid")
    if err != nil {
        return
    }
    uuidString, _ := strings.CutPrefix(uuidCookie.String(), "uuid=")
    if uuidString == "" {
        return
    }

    value, got := cache.UUID.Get(uuidString)
    if got {
        return value
    }
    return
}


