package main

import (
	"html/template"
	"net/http"
	"time"
)

type LoginForm struct {
    Username string
    Password string
    Remember bool
    admin bool
}

type UploadForm struct {
    User string
    Messages []string
    Overwrite bool
}

var tmpl *template.Template

func init() {
    tmpl = template.Must(template.ParseFiles(
        "templates/index.html",
        "templates/login.html",
        "templates/upload.html",
        ))
}

// The root of the website
func MainHandler(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "base", struct{Logged bool}{FromCookie(r).Name != ""})
}

// /login/
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
        tmpl.ExecuteTemplate(w, "login", nil)
        return
    }

    var form LoginForm
    form.Username = r.PostFormValue("username")
    form.Password = r.PostFormValue("password")
    if r.PostFormValue("remember") == "on" {
        form.Remember = true
    }

    // check credentials
    if form.Username == "" {
        //tmpl.ExecuteTemplate(w, "login", nil)
        return
    }

    ip := getClientIP(r)
    if _, got := Limited.Get(ip); got {
        w.Write([]byte("<p>Please wait before trying again"))
        return
    }
    // TODO: countdown ?

    if CheckCredentials(&form, &w) {
        // TODO Maybe pass username ?
        w.Header().Set("HX-Retarget", "#main-form")
        tmpl.ExecuteTemplate(w, "upload", nil)
        return
    } else {
        Limited.Set(ip, true, time.Duration(Conf.Login_ttl))
    }

    w.Write([]byte("<p>Invalid Username/Password"))
    //tmpl.ExecuteTemplate(w, "login", fields)
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

// /upload/
func UploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") != "true" {
        return
    }

    if user := FromCookie(r); user.Name == "" {
        // TODO: execute main with HX-Boost and message
        w.Header().Set("HX-Retarget", "#main-form")
        tmpl.ExecuteTemplate(w, "login", nil)
        w.Header().Set("HX-Retarget", "#messages")
        w.Write([]byte("<p>Login Expired"))
        return
    }

    var form UploadForm
    if r.PostFormValue("overwrite") == "on" {
        form.Overwrite = true
    }

    // handle files
    if files := r.MultipartForm.File["file"]; files != nil {
        HandleFiles(&form, files)
    } else {
        form.Messages = append(form.Messages, "No file chosen")
    }

    for _, msg := range form.Messages {
        w.Write([]byte("<p>" + msg))
    }
    //Tmpl.upload.Execute(w, form)
}

// TODO: CPanel handler
func MenuHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") != "true" {
        return
    }

    user := FromCookie(r)
    if user.Name == "" {
        return
    }
}
