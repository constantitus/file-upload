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

var UUID cache.Cache[string, Data]

func init() {
    UUID = cache.NewCache[string, Data]()
    go func() {
        for {
            time.Sleep(time.Minute * 1)
            UUID.DeleteExpired()
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

    UUID.Set(uuidString, Data{}, time.Duration(1))
    return
}

