package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var Tmpl struct {
    header *template.Template
    footer *template.Template
    login *template.Template
    upload *template.Template
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
    Tmpl.header.Execute(w, nil)
    if user := FromCookie(r); user != "" {
        Tmpl.upload.Execute(w, UploadForm{User: user})
    } else {
        Tmpl.login.Execute(w, nil)
    }
    Tmpl.footer.Execute(w, nil)
}


func init() {
    Tmpl.header = template.Must(template.ParseFiles("templates/header.html"))
    Tmpl.footer = template.Must(template.ParseFiles("templates/footer.html"))
    Tmpl.login = template.Must(template.ParseFiles("templates/login.html"))
    Tmpl.upload = template.Must(template.ParseFiles("templates/upload.html"))
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
