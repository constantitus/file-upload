package main

import (
	"fmt"
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
        ))
}

// The root of the website
func MainHandler(w http.ResponseWriter, r *http.Request) {
    // TODO: Pass user data to upload template
    data := FromCookie(r)
    args := struct{
        Logged bool
        Files []DirEntry // temporary solution
    }{}
    if data.Name != "" {
        args.Logged = true
        args.Files = ReadUserDir(data.Name)
    }
    w.Write([]byte(Style))
    tmpl.ExecuteTemplate(w, "base", args)
}


type LoginData struct {
    Username string
    Password string
    Remember bool
    admin bool
    Files []DirEntry
}
// /login/
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") != "true" {
        return
    }

    // Logout
    if r.PostFormValue("logout") == "true" {
        ClearFromCache(r)
        cookie := &http.Cookie{
            Name: "uuid",
            Path: "/",
            Expires: time.Now(),
        }
        http.SetCookie(w, cookie)
        w.Header().Set("HX-Retarget", "#main")
        tmpl.ExecuteTemplate(w, "login", nil)
        return
    }

    var form LoginData
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
            w.Header().Set("HX-Retarget", "#main")
            // TODO: Pass userdata
            tmpl.ExecuteTemplate(w, "upload", struct{
                Files []DirEntry
            }{
                ReadUserDir(form.Username), // temporary solution
            })
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


type UploadData struct {
    User string
    Messages []string
    Overwrite bool
}
// /upload/
func UploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") != "true" {
        return
    }

    var form UploadData
    user := FromCookie(r)
    if user.Name == "" {
        w.Header().Set("HX-Retarget", "#main")
        tmpl.ExecuteTemplate(w, "login", nil)
        return
    }
    form.User = user.Name

    if r.PostFormValue("overwrite") == "on" {
        form.Overwrite = true
    }

    // handle files
    var entries []DirEntry
    if files := r.MultipartForm.File["file"]; files != nil {
        HandleFiles(&form, files, &entries)
    } else {
        form.Messages = append(form.Messages, "No file chosen")
    }

    // w.Write([]byte(`<div id="">`))
    // w.Write([]byte(`</div>`))
    // Files (HAS TO GO FIRST)
    // w.Write([]byte(`<tbody hx-swap-oob="beforeend:#directory">`))
    tmpl.ExecuteTemplate(w, "file", struct{Files []DirEntry}{entries})
    // w.Write([]byte(`</tbody>`))
    // Messages
    for _, msg := range form.Messages {
        w.Write([]byte(`
    <p>` + msg))
    }
}


// /files/
func FileHandler(w http.ResponseWriter, r *http.Request) {
    // handle the file download
    if r.Method == "GET" {
        query := r.URL.Query()
        user, got := UUID.Get(query.Get("uuid"))
        file := query.Get("download")
        if !got || file == "" { return }
        w.Header().Set("Content-Disposition", "attachment; filename=" + file)
        w.Header().Set("Content-Type", "application/octet-stream")
        http.ServeFile(w, r, Config.StoragePath + "/" + user.Name + "/" + file)
        return
    }

    if r.Header.Get("HX-Request") != "true" {
        return
    }

    user := FromCookie(r)
    if user.Name == "" {
        return
    }

    params := []string{"download", "delete", "rename"}
    for _, param := range params {
        val := r.PostFormValue(param)
        if val == "" {
            continue
        }
        switch param {
            case "download":
            uuid, err := r.Cookie("uuid")
            if err != nil {
                return
            }
            w.Header().Set(
                "HX-Redirect",
                fmt.Sprintf("/files?uuid=%s&download=%s", uuid.Value, val),
                )
            return
            case "delete":
            // TODO: handle delete
            case "rename":
            tmpl.ExecuteTemplate(w, "rename", nil)
            newname := r.PostFormValue("newname")
            if newname == "" {
            }
        }
    }
}

