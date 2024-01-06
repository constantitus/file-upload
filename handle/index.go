package handle

import (
	"html/template"
	"net/http"
	"nullprogram.com/x/uuid"
)

var tmpl *template.Template
var gen *uuid.Gen

func init() {
    gen = uuid.NewGen()

    tmpl = template.Must(template.ParseFiles(
        "templates/index.html",
        "templates/login.html",
        "templates/upload.html",
        "templates/table.html",
        ))
}

// The root of the website
func IndexHandler(w http.ResponseWriter, r *http.Request) {
    // TODO: Pass user data to upload template
    data := fromCookie(r)
    args := struct{
        Logged bool
        Files []dirEntry // temporary solution
    }{}
    if data.Name != "" {
        args.Logged = true
        args.Files = readUserDir(data.Name)
    }
    tmpl.ExecuteTemplate(w, "base", args)
}
