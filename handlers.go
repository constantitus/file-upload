package main

import (
	"html/template"
	"net/http"
	"time"
)

var Tmpl struct {
    header *template.Template
    footer *template.Template
    login *template.Template
    upload *template.Template
}

type LoginForm struct {
    Username string
    Password string
    Remember bool
    Message string
    admin bool
}

type UploadForm struct {
    User string
    Message []string
    Overwrite bool
}

func init() {
    Tmpl.header = template.Must(template.ParseFiles("templates/header.html"))
    Tmpl.footer = template.Must(template.ParseFiles("templates/footer.html"))
    Tmpl.login = template.Must(template.ParseFiles("templates/login.html"))
    Tmpl.upload = template.Must(template.ParseFiles("templates/upload.html"))
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") != "true" {
        return
    }

    // Delete cookies on logout
    if r.PostFormValue("logout") == "true" {
        cookie := &http.Cookie{
            Name: "uuid",
            Path: "/",
            Expires: time.Now(),
        }
        http.SetCookie(w, cookie)
        Tmpl.login.Execute(w, nil)
        return
    }

    var fields LoginForm
    fields.Username = r.PostFormValue("username")
    fields.Password = r.PostFormValue("password")
    if r.PostFormValue("remember") == "on" {
        fields.Remember = true
    }

    // check credentials
    if fields.Username == "" {
        Tmpl.login.Execute(w, LoginForm{Remember: fields.Remember})
        return
    }

    if CheckCredentials(&fields, &w) {
        Tmpl.upload.Execute(w, UploadForm{User: fields.Username})
        return
    }

    Tmpl.login.Execute(w, fields)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") != "true" {
        return
    }

    if user := FromCookie(r); user == "" {
        Tmpl.login.Execute(w, LoginForm{Message: "Login Expired"})
        return
    }

    var form UploadForm
    if r.PostFormValue("overwrite") == "on" {
        form.Overwrite = true
    }

    // handle files
    HandleFiles(&form, r.MultipartForm.File["file"])
}
