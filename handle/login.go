package handle

import (
	"net/http"
	"time"

	"main/cache"
	"main/limits"
)

type loginData struct {
    Username string
    Password string
    Remember bool
    admin bool
    Files []dirEntry
}
// /login/
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") != "true" {
        return
    }

    // Logout
    if r.PostFormValue("logout") == "true" {
        cache.RemoveEntry(r)
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

    var form loginData
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

    if !limits.CheckLimit(ip) {
        w.Write([]byte("<p>Please wait before trying again"))
        return
    } else {
        if checkCreds(&form, &w) {
            w.Header().Set("HX-Retarget", "#main")
            // TODO: Pass userdata
            tmpl.ExecuteTemplate(w, "upload", struct{
                Files []dirEntry
            }{
                readUserDir(form.Username), // temporary solution
            })
            return
        } else {
            w.Write([]byte("<p>Invalid Username/Password"))
            return
        }
    }
}

