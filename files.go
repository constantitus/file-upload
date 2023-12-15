package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
)

type File struct {
    name string
    content *[]byte
}

func HandleFiles(form *UploadForm, headers []*multipart.FileHeader) {
    // Parse files
    var files []File
    for _, f := range headers {
        // TODO: limit f.Size()
        file, err := f.Open()
        if err != nil {
            log.Println("Upload error:", err.Error())
            continue
        }
        defer file.Close()

        /* // This is from io.ReadAll()
        var read int64
        buf := make([]byte, 0, 512)
        for {
            size, err := file.Read(buf[len(buf):cap(buf)])
            buf = buf[:len(buf)+size]
            if err == io.EOF {
                break
            }
            if len(buf) == cap(buf) {
                buf = append(buf, 0)[:len(buf)]
            }
            read = read + int64(size)
        } */
        buf, err := io.ReadAll(file)
        if err != nil {
            form.Messages = append(form.Messages, err.Error())
            break
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
            form.Messages = append(form.Messages, "No file selected")
            continue
        }
        out, err := os.OpenFile(
            Conf.StoragePath + "/" + form.User + "/" + file.name,
            os.O_RDWR|os.O_CREATE,
            0644)
        defer out.Close()
        if err != nil {
            form.Messages = append(form.Messages, err.Error())
            continue
        }
        out_stat, err := out.Stat()
        if err != nil {
            log.Println("our.Stat()", err)
        }
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
        form.Messages = append(form.Messages, msg)
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
