package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type File struct {
    name string
    content *[]byte
}

type UploadForm struct {
    Message []string
    Overwrite bool
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") != "true" {
        return
    }
    // check creds
    if !CheckCookies(r) {
        Tmpl.login.Execute(w, LoginForm{Message: "Invalid Credentials"})
        return
    }

    var form UploadForm
    if r.PostFormValue("overwrite") == "on" {
        form.Overwrite = true
    }

    // handle files
    var files []File
    fileHeaders := r.MultipartForm.File["file"]
    for _, f := range fileHeaders {
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

    // handle writing
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
                msg += Write(&file, *out)
            } else {
                msg += "File already exists: " + file.name
            }
        } else {
            msg += Write(&file, *out)
        }
        form.Message = append(form.Message, msg)
    }

    // TODO: Parse <p>'s in form.Message
    Tmpl.upload.Execute(w, form)
}

func Write(file *File, out os.File) string {
    out_size, err := out.Write(*file.content)
    if err != nil {
        return err.Error()
    } else {
        return fmt.Sprintf("written %s (%d bytes)", file.name, out_size)
    }
}


/* func UploadHandlerOld(w http.ResponseWriter, r *http.Request) {
    reader, err := r.MultipartReader()
    if err != nil {
        fmt.Println("Error: ", err.Error())
        return
    }

    form, err := reader.ReadForm(1000 << 20)
    if err != nil {
        fmt.Println("Error: ", err.Error())
        return
    }

    var files []File
    for _, headers := range form.File {
        for _, header := range headers {
            file, err := header.Open()
            if err != nil {
                fmt.Println("Error: ", err.Error())
                continue
            }
            defer file.Close()

            buf := bytes.NewBuffer(nil)
            _, err = io.Copy(buf, file)
            if err != nil {
                fmt.Println("Error: ", err.Error())
                return
            }

            //files = append(files, File{header.Filename, &buf.Bytes()})
        }
    }
} */
