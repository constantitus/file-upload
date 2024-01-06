package handle

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

    "main/config"
)

// parses the files to be uploaded by the user
func parseUpload(form *uploadData, headers []*multipart.FileHeader) {
    // Parse files
    for _, header := range headers {
        filename := header.Filename

        if filename == "" {
            form.Messages = append(form.Messages, "No file selected")
            continue
        }

        if header.Size > config.FilesizeMax {
            form.Messages = append(form.Messages, "file too large: " + filename)
            continue
        }

        // this is in case the user doesn't have a dir yet
        os.MkdirAll(config.StoragePath + "/" + form.User, os.ModePerm)

        out, err := os.OpenFile(
            config.StoragePath +  form.User + "/" + filename,
            os.O_RDWR|os.O_CREATE,
            0644)
        defer out.Close()
        if err != nil {
            log.Println(err)
            form.Messages = append(form.Messages, "internal error " + filename)
            continue
        }

        out_stat, err := out.Stat()
        if err != nil { /* why would this even error */ }

        if out_stat.Size() == 0 || form.Overwrite {
            written, err := storeFile(header, *out)
            if err != nil {
                log.Println(err)
                form.Messages = append(form.Messages, "internal error " + filename)
                continue
            }

            var msg string
            if out_stat.Size() > 0 && form.Overwrite {
                msg = "over" // "over" + "written", duh
            }
            msg += fmt.Sprintf(
                "written %s (%s)",
                filename,
                sizeItoa(written),
                )
            form.Messages = append(form.Messages, msg)
            continue
        }
        form.Messages = append(
            form.Messages, "File already exists: " + filename)
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
    }
    return out_size, nil
}


type dirEntry struct {
    Name string
    Size string
    Date string
}
// Returns the files in the user directory
func readUserDir(username string) []dirEntry {
    var entries []dirEntry
    dir, err := os.ReadDir(config.StoragePath + "/" + username)
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

        entries = append(entries, dirEntry{
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


func tryRename(user string, old string, newname string) (bool, string) {
    if sanitize(old) {
        return false, `Illegal character`
    }
    if sanitize(newname) {
        return false, "Illegal character"
    }

    err := os.Rename(
        config.StoragePath + "/" + user + "/" + old,
        config.StoragePath + "/" + user + "/" + newname,
        )
    if err != nil {
        return false, "Internal Error"
    }
    return true, "renamed " + old + " to " + newname
}

func tryRemove(user string, filename string) (bool, string) {
    if sanitize(filename) {
        return false, "Illegal Character"
    }

    err := os.Remove( config.StoragePath + "/" + user + "/" + filename,)
    if err != nil {
        return false, "Internal Error"
    }
    return true, "deleted " + filename

}

func sanitize(filename string) bool {
    return strings.ContainsRune(filename, '/')
}
