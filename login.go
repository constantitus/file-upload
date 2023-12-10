package main

import (
	"net/http"
)

type LoginForm struct {
    Username string
    Password string
    Remember bool
    Message string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") != "true" {
        return
    }

    // Delete cookies on logout
    if r.PostFormValue("logout") == "true" {
        cookies := []string{"username", "password"}
        for _, cookie := range cookies {
            cookie := &http.Cookie{
                Name: cookie,
                Value: "",
                Path: "/",
            }
            http.SetCookie(w, cookie)
        }
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
        Tmpl.upload.Execute(w, nil)
        return
    }

    Tmpl.login.Execute(w, fields)
}



/* func LoginHandlerV1(w http.ResponseWriter, r *http.Request) {
    // I am a dumbass.
    // I always manage to get shit done in the most complicated
    // way without even thinking that ther might be an easier way
    reader, err := r.MultipartReader()
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    form, err := reader.ReadForm(1000 << 20)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    var fields Fields
    for name, value:= range form.Value {
        switch name {
        case "username":
            fields.username = check(value)
        case "password":
            fields.password = check(value)
        case "remember":
            if check(value) == "on" {
                fields.remember = true
            }
        }
    }
    // fmt.Println(fields)
}

func check(in []string) string {
    if len(in) > 0 {
        return in[0]
    }
    return ""
} */
