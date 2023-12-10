package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"github.com/jellydator/ttlcache/v3"
)

var Tmpl struct {
    header *template.Template
    footer *template.Template
    login *template.Template
    upload *template.Template
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
    Tmpl.header.Execute(w, nil)
    if CheckCookies(r) {
        Tmpl.upload.Execute(w, nil)
    } else {
        Tmpl.login.Execute(w, nil)
    }
    Tmpl.footer.Execute(w, nil)
}

var cache *ttlcache.Cache[string, string]

func init() {
    Tmpl.header = template.Must(template.ParseFiles("templates/header.html"))
    Tmpl.footer = template.Must(template.ParseFiles("templates/footer.html"))
    Tmpl.login = template.Must(template.ParseFiles("templates/login.html"))
    Tmpl.upload = template.Must(template.ParseFiles("templates/upload.html"))

    cache = ttlcache.New[string, string](
        ttlcache.WithTTL[string, string](30 * time.Minute),
    )
    go cache.Start()
}

func main() {
    fmt.Println("Server running...")
    http.HandleFunc("/", MainHandler)
    http.HandleFunc("/upload/", UploadHandler)
    http.HandleFunc("/login/", LoginHandler)
    server := http.Server{
        Addr:         ":8080",
        Handler:      nil,
        // ReadTimeout:  10000,
        // WriteTimeout: 10000,
    }

    log.Fatal(server.ListenAndServe())
}
