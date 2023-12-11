package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type File struct {
    name string
    content *[]byte
}

type UploadForm struct {
    User string
    Message []string
    Overwrite bool
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") != "true" {
        return
    }

    if user := FromCookie(r); user == "" {
        Tmpl.login.Execute(w, LoginForm{Message: "Invalid Credentials"})
        return
    }

    var form UploadForm
    if r.PostFormValue("overwrite") == "on" {
        form.Overwrite = true
    }

    // handle files
    handleFiles(&form, r.MultipartForm.File["file"])

}

func handleFiles(form *UploadForm, headers []*multipart.FileHeader) {
    // Parse files
    var files []File
    for _, f := range headers {
        file, err := f.Open()
        if err != nil {
            fmt.Println(err.Error())
            continue
        }
        defer file.Close()
        buf, err := io.ReadAll(file)
        if err != nil {
            form.Message = append(form.Message, err.Error())
        }

        files = append(files, File{
            name: f.Filename,
            content: &buf,
        })
    }

    // Write files
    for _, file := range files {
        var msg string
        if file.name == "" {
            msg += "No file selected"
            continue
        }
        out, err := os.OpenFile(
            file.name, // TODO: handle directories
            os.O_RDWR|os.O_CREATE,
            0644)
        defer out.Close()
        if err != nil {
            msg += err.Error()
            continue
        }
        out_stat, _ := out.Stat()
        exists := out_stat.Size() != 0
        if exists {
            if form.Overwrite {
                msg += "over" // overwritten
                msg += writeFile(&file, *out)
            } else {
                msg += "File already exists: " + file.name
            }
        } else {
            msg += writeFile(&file, *out)
        }
        form.Message = append(form.Message, msg)
    }
}

func writeFile(file *File, out os.File) string {
    out_size, err := out.Write(*file.content)
    if err != nil {
        return err.Error()
    } else {
        return fmt.Sprintf("written %s (%d bytes)", file.name, out_size)
    }
}
