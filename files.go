package main

import (
	"fmt"
	"io"
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

    // THIS DOESN T HAPPEN
    // Write files
    var msg string
    for _, file := range files {
        if file.name == "" {
            msg += "No file selected"
            continue
        }
        out, err := os.OpenFile(
            Conf.Storage_path + "/" + file.name,
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
    }
    form.Message = append(form.Message, msg)
}

func writeFile(file *File, out os.File) string {
    out_size, err := out.Write(*file.content)
    if err != nil {
        return err.Error()
    } else {
        return fmt.Sprintf("written %s (%d bytes)", file.name, out_size)
    }
}
