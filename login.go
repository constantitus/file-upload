package main

import (
	"net/http"
	"time"
)

type LoginForm struct {
    Username string
    Password string
    Remember bool
    Message string
    admin bool
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
