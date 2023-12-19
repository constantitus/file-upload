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

func HandleFiles(form *UploadData, headers []*multipart.FileHeader, entries *[]DirEntry) {
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
        os.MkdirAll(Config.StoragePath + "/" + form.User, os.ModePerm)
        out, err := os.OpenFile(
            Config.StoragePath + "/" + form.User + "/" + file.name,
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
                msg += writeFile(&file, *out, entries)
            } else {
                msg += "File already exists: " + file.name
            }
        } else {
            msg += writeFile(&file, *out, entries)
        }
        form.Messages = append(form.Messages, msg)
    }
}

func writeFile(file *File, out os.File, entries *[]DirEntry) string {
    out_size, err := out.Write(*file.content)
    if err != nil {
        return err.Error()
    } else {
        *entries = append(*entries, DirEntry{Name: file.name, Size: int64(out_size)})
        return fmt.Sprintf("written %s (%d bytes)", file.name, out_size)
    }
}


type DirEntry struct {
    Name string
    Size int64
}
func ReadUserDir(username string) []DirEntry {
    var entries []DirEntry
    dir, err := os.ReadDir(Config.StoragePath + "/" + username)
    if err != nil {
        log.Println(err.Error())
        return entries
    }

    for _, entry := range dir {
        if entry.IsDir() { continue }
        info, err := entry.Info()
        if err != nil {
            fmt.Println(err.Error())
            continue
        }
        entries = append(entries, DirEntry{
            Name: entry.Name(),
            Size: info.Size(),
        })
    }
    return entries
}
