package handle

import (
    "net/http"
	"fmt"

	"main/cache"
	"main/config"
)
// /files/
func OptionsHandler(w http.ResponseWriter, r *http.Request) {
    // handle the file download
    if r.Method == "GET" {
        query := r.URL.Query()
        user, got := cache.UUID.Get(query.Get("uuid"))
        file := query.Get("download")
        if !got || file == "" { return }
        w.Header().Set("Content-Disposition", "attachment; filename=" + file)
        w.Header().Set("Content-Type", "application/octet-stream")
        http.ServeFile(w, r, config.StoragePath + "/" + user.Name + "/" + file)
        return
    }

    if r.Header.Get("HX-Request") != "true" {
        return
    }
    user := fromCookie(r)
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
            success, msg := tryRemove(user.Name, entry)
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
            success, msg := tryRename(user.Name, entry, reply.NewName)
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
        "multi:#file-browser:outerHTML,#extra:innerHTML,#messages:innerHTML",
        )

    // update table
    entries := readUserDir(username)
    tmpl.ExecuteTemplate(w, "file-table", struct{Files []dirEntry}{entries})

    // close prompt
    w.Write([]byte(`<div id="extra"></div>`))

    // print messages
    tmpl.ExecuteTemplate(w, "message", msg)
}
