package handle

import (
	"net/http"
)

type uploadData struct {
    User string
    Messages []string
    Overwrite bool
}
// /upload/
func Upload(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") != "true" {
        return
    }

    var form uploadData
    user := fromCookie(r)
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
        parseUpload(&form, files)
    } else {
        form.Messages = append(form.Messages, "No file chosen")
    }

    w.Header().Set("HX-Reswap", "multi:#file-browser:outerHTML,#messages")

    // Update files - We're refreshing the whole table.
    // While we could add new elements, htmx has nothing that can allow us to
    // modify an existing table element. It'd be too much of a hassle anyway.
    entries := readUserDir(user.Name)
    tmpl.ExecuteTemplate(w, "file-table", struct{Files []dirEntry}{entries})

    // Print messages
    tmpl.ExecuteTemplate(w, "messages", form.Messages)
}

