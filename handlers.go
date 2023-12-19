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

    // Deletes UUID from memcache and cookies on logout
    if r.PostFormValue("logout") == "true" {
        // clear the cache if the login is valid
        ClearFromCache(r)
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
            w.Header().Set("HX-Swap", "outerHTML")
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
    var entries []DirEntry
    if files := r.MultipartForm.File["file"]; files != nil {
        HandleFiles(&form, files, &entries)
    } else {
        form.Messages = append(form.Messages, "No file chosen")
    }

    // Files (HAS TO GO FIRST)
    tmpl.ExecuteTemplate(w, "file", struct{Entries []DirEntry}{entries})
    /* w.Write([]byte(`<tr class="file">
        <td>pls</td>
        <td>work</td>
        <td><button>download</button></td>
        <td><button>delete</button></td>
        <td><button>rename</button></td>
    </tr>`)) */
    // Messages
    w.Write([]byte(`<div id="messages">`))
    for _, msg := range form.Messages {
        w.Write([]byte(`
    <p>` + msg))
    }
    w.Write([]byte(`</div>`))
}



type FbReply struct {
    Entries []DirEntry
}
// Test handler
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

func FileDownload() {}

func FileDelete() {}

func FileReneme() {}
