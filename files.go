package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
)


func HandleFiles(form *UploadData, headers []*multipart.FileHeader) {
    // Parse files
    for _, header := range headers {
        // TODO: limit f.Size()

        if header.Filename == "" {
            form.Messages = append(form.Messages, "No file selected")
            continue
        }

        // this is in case the user doesn't have a dir yet
        os.MkdirAll(Config.StoragePath + "/" + form.User, os.ModePerm)

        out, err := os.OpenFile(
            Config.StoragePath + "/" + form.User + "/" + header.Filename,
            os.O_RDWR|os.O_CREATE,
            0644)
        defer out.Close()
        if err != nil {
            form.Messages = append(form.Messages, err.Error())
            continue
        }

        // 
        out_stat, err := out.Stat()
        if err != nil { /* why would this even error */ }

        if out_stat.Size() == 0 || form.Overwrite {
            written, err := storeFile(header, *out)
            if err != nil {
                form.Messages = append(form.Messages, err.Error())
                continue
            }

            var msg string
            if out_stat.Size() > 0 && form.Overwrite {
                msg = "over" // "overwritten"
            }
            msg += fmt.Sprintf(
                "written %s (%s)",
                header.Filename,
                sizeItoa(written),
                )
            form.Messages = append(form.Messages, msg)
            continue
        }
        form.Messages = append(
            form.Messages, "File already exists: " + header.Filename)
    }
}
// The actual reading and writing
func storeFile(header *multipart.FileHeader, out os.File) (int64, error) {
    file, err := header.Open()
    defer file.Close()
    if err != nil {
        return 0, err
    }

    out_size, err := io.Copy(&out, file)
    if err != nil {
        return 0, err
    } else {
        return out_size, nil
    }
}


type DirEntry struct {
    Name string
    Size string
    Date string
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
            Size: sizeItoa(info.Size()),
            Date: parseDate(info.ModTime()),
        })
    }
    return entries
}

func sizeItoa(size int64) (out string) {
    in := int(size)
    if i := (in >> 30); i > 0 {
        return strconv.Itoa(i) + "GB"
    }
    if i := (in >> 20); i > 0 {
        return strconv.Itoa(i) + "MB"
    }
    if i := (in >> 10); i > 0 {
        return strconv.Itoa(i) + "KB"
    } else {
        return strconv.Itoa(in) + "B"
    }
}

func parseDate(d time.Time) string {
    if since := time.Since(d); since < 24 * time.Hour {
        since = since.Round(time.Second)
        if since < 1 { return "now" }
        if since < time.Minute {
            return fmt.Sprintf("%02d sec ago", since / time.Second)
        }
        if since < time.Hour {
            return fmt.Sprintf("%02d min ago", since / time.Minute)
        }
        return fmt.Sprintf("%02d hr ago", since / time.Hour)
    }
    mon := d.Month().String()[:3]
    if d.Year() == time.Now().Year() {
        return fmt.Sprintf("%d %s", d.Day(), mon)
    } else {
        return fmt.Sprintf("%s %d", mon, d.Year())
    }
}


func TryRename(user string, old string, newname string) (bool, string) {
    if sanitize(old) {
        return false, `Illegal character`
    }
    if sanitize(newname) {
        return false, "Illegal character"
    }

    err := os.Rename(
        Config.StoragePath + "/" + user + "/" + old,
        Config.StoragePath + "/" + user + "/" + newname,
        )
    if err != nil {
        return false, "Internal Error"
    }
    return true, "renamed " + old + " to " + newname
}

func TryRemove(user string, filename string) (bool, string) {
    if sanitize(filename) {
        return false, "Illegal Character"
    }

    err := os.Remove( Config.StoragePath + "/" + user + "/" + filename,)
    if err != nil {
        return false, "Internal Error"
    }
    return true, "deleted " + filename

}

func sanitize(filename string) bool {
    return strings.ContainsRune(filename, '/')
}
