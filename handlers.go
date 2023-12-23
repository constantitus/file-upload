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
        "templates/table.html",
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
    if files := r.MultipartForm.File["file"]; files != nil {
        HandleFiles(&form, files)
    } else {
        form.Messages = append(form.Messages, "No file chosen")
    }

    w.Header().Set("HX-Reswap", "multi:#file-browser:outerHTML,#messages")

    // Update files - We're refreshing the whole table.
    // While we could add new elements, htmx has nothing that can allow us to
    // modify an existing table element. It'd be too much of a hassle anyway.
    entries := ReadUserDir(user.Name)
    tmpl.ExecuteTemplate(w, "file-table", struct{Files []DirEntry}{entries})

    // Print messages
    w.Write([]byte(`<div id="messages">`))
    for _, msg := range form.Messages {
        w.Write([]byte(`
    <p>` + msg))
    }
    w.Write([]byte(`</div>`))
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

    option := r.PostFormValue("option")
    entry := r.PostFormValue("entry")
    if entry == "" || option == "" {
        return // neither should be empty
    }


    reply := struct {
        Name string;
        Message string;
        NewName string; // keep value for rename
    }{}
    switch option {
    case "download":
        uuid, err := r.Cookie("uuid")
        if err != nil {
            return
        }
        w.Header().Set(
            "HX-Redirect",
            fmt.Sprintf("/files?uuid=%s&download=%s", uuid.Value, entry),
            )
        return

    case "delete":
        if delete := r.PostFormValue("delete"); delete == "yes" {
            // delete
            success, msg := TryRemove(user.Name, entry)
            if success {
                onSuccess(w, msg, user.Name)
                return
            }
            reply.Message = msg
        }
        reply.Name = entry
        tmpl.ExecuteTemplate(w, "delete", struct{Name string}{entry})

    case "rename":
        reply.NewName = r.PostFormValue("newname")
        if reply.NewName != "" {
            success, msg := TryRename(user.Name, entry, reply.NewName)
            if success {
                onSuccess(w, msg, user.Name)
                return
            }
            reply.Message = msg
        } else {
            reply.NewName = entry
        }
        reply.Name = entry
        tmpl.ExecuteTemplate(w, "rename", reply)
    }
}

func onSuccess(w http.ResponseWriter, msg string, username string) {
    w.Header().Set("HX-Reswap",
        "multi:#file-browser:outerHTML,#pop-window:delete,#messages:innerHTML",
        )

    // update table
    entries := ReadUserDir(username)
    tmpl.ExecuteTemplate(w, "file-table", struct{Files []DirEntry}{entries})

    // close prompt
    w.Write([]byte(`<div id="pop-window"></div>`))

    // print messages
    tmpl.ExecuteTemplate(w, "message", msg)
}
