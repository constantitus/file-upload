package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
)

type Cookies struct {
    user http.Cookie
    pass http.Cookie
}

func InitCookies(fields *LoginForm, pass string) (t Cookies) {
    var expires time.Time
    if fields.Remember {
        expires = time.Now().Add(time.Hour * 24)
    }
    t.user = http.Cookie{
        Name: "username",
        Value: fields.Username,
        Expires: expires,
        Path: "/",
    }
    t.pass = http.Cookie{
        Name: "password",
        Value: pass,
        Expires: expires,
        Path: "/",
    }
    return
}

func CheckCredentials(fields *LoginForm, w *http.ResponseWriter) bool {
    /* if Users[fields.Username] == "" {
        fields.Message = "Invalid credentials"
        return false
    } */

    tmp := sha256.New()
    tmp.Write([]byte(fields.Password))
    hashed_pass := hex.EncodeToString(tmp.Sum(nil))

    if hashed_pass == CheckUserDB(fields.Username) {
        cookies := InitCookies(fields, hashed_pass)
        http.SetCookie(*w, &cookies.user)
        http.SetCookie(*w, &cookies.pass)

        return true
    }
    return false
}

func CheckCookies(r *http.Request) bool {
    user_cookie, err := r.Cookie("username")
    if err != nil {
        return false
    }
    pass_cookie, err := r.Cookie("password")
    if err != nil {
        return false
    }
    user, _ := strings.CutPrefix(user_cookie.String(), "username=")
    pass, _ := strings.CutPrefix(pass_cookie.String(), "password=")
    if user == "" || pass == "" {
        return false
    }

    return CheckUserDB(user) == pass
}
