package cache

import (
	"net/http"
	"strings"
	"time"

    "github.com/go-pkgz/expirable-cache/v2"
)

type Data struct {
    Name string
    Rank bool
}

var uuid cache.Cache[string, Data]

func init() {
    uuid = cache.NewCache[string, Data]()
    go func() {
        for {
            time.Sleep(time.Minute * 1)
            uuid.DeleteExpired()
        }
    }()
}

// on logout
func RemoveEntry(r *http.Request) {
    uuidCookie, err := r.Cookie("uuid")
    if err != nil {
        return
    }
    uuidString, _ := strings.CutPrefix(uuidCookie.String(), "uuid=")
    if uuidString == "" {
        return
    }

    uuid.Set(uuidString, Data{}, time.Duration(1))
    return
}

// Exposed methods in case I want to change the cache

func Keys() []string {
    return uuid.Keys()
}

func Get(key string) (Data, bool) {
    return uuid.Get(key)
}

func GetExp(key string) (time.Time, bool) {
    return uuid.GetExp(key)
}

func Set(key string, value Data, ttl time.Duration) {
    uuid.Set(key, value, ttl)
}

func DeleteExpired() {
    uuid.DeleteExpired()
}
