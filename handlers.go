package main

import (
	"html/template"
	"net/http"
	"time"
)

var tmpl *template.Template

func init() {
    tmpl = template.Must(template.ParseFiles(
        "templates/index.html",
        "templates/login.html",
        "templates/upload.html",
        "templates/menu.html",
        ))
}

// The root of the website
func MainHandler(w http.ResponseWriter, r *http.Request) {
    args := struct{
        Logged bool
    }{
        FromCookie(r).Name != "",
    }
    w.Write([]byte(Style))
    tmpl.ExecuteTemplate(w, "base", args)

}


type LoginForm struct {
    Username string
    Password string
    Remember bool
    admin bool
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
        w.Header().Set("HX-Swap", "outerHTML")
        w.Header().Set("HX-Retarget", "#main")
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
        return
    }

    ip := getClientIP(r)

    if !CheckLimit(ip) {
        w.Write([]byte("<p>Please wait before trying again"))
        return
    } else {
        if CheckCredentials(&form, &w) {
            // TODO Maybe pass username ?
            w.Header().Set("HX-Swap", "outerHTML")
            w.Header().Set("HX-Retarget", "#main")
            tmpl.ExecuteTemplate(w, "upload", nil)
            return
        } else {
            w.Write([]byte("<p>Invalid Username/Password"))
            return
        }
    }
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


type UploadForm struct {
    User string
    Messages []string
    Overwrite bool
}
// /upload/
func UploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") != "true" {
        return
    }

    var form UploadForm
    user := FromCookie(r)
    if user.Name == "" {
        w.Header().Set("HX-Swap", "outerHTML")
        w.Header().Set("HX-Retarget", "#main")
        tmpl.ExecuteTemplate(w, "login", nil)
        return
    }
    form.User = user.Name

    if r.PostFormValue("overwrite") == "on" {
        form.Overwrite = true
    }

    // handle files
    if files := r.MultipartForm.File["file"]; files != nil {
        HandleFiles(&form, files)
    } else {
        form.Messages = append(form.Messages, "No file chosen")
    }

    // w.Header().Set("HX-Retarget", "#messages")
    for _, msg := range form.Messages {
        w.Write([]byte("<p>" + msg))
    }
}


type FbReply struct {
    Entries []DirEntry
}

func FileBrowserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") != "true" {
        return
    }

    user := FromCookie(r)
    if user.Name == "" {
        return
    }

    reply := FbReply{Entries: ReadUserDir(user.Name)}


    w.Header().Set("HX-Swap", "outerHTML")
    w.Header().Set("HX-Retarget", "#main")
    tmpl.ExecuteTemplate(w, "menu", reply)
}
