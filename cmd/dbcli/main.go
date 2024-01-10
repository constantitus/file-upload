package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"main/db"
	"os"
	"strconv"
)

const (
    usage = `usage: db-cli [options] [arguments]
    list    - List users and their rank
    add     - Add user to the database
    rm      - Remove user from the database`

    addUsage = `usage: db-cli add [username] [password] [admin]
    (admin is boolean)`

    rmUsage = `usage: db-cli rm [username]`
)

func main() {
    argv := os.Args
    argc := len(argv)

    if argc <= 1 {
        fmt.Println(usage)
        return
    }

    switch argv[1] {
    case "add":
        if argc == 5 {
            addUser(argv[2], argv[3], argv[4])
            return
        }
        fmt.Println(addUsage)
    case "rm":
        if argc == 3 {
            rmUser(argv[2])
            return
        }
        fmt.Println(rmUsage)
    case "list":
        list()
        return
    default:
        fmt.Println("unknown option: ", argv[1])
        return
    }

}


func addUser(user string, pass string, admin string) {
    if err := db.Read(); err != nil {
        log.Fatalln(err)
    }

    tmp := sha256.New()
    tmp.Write([]byte(pass))
    hash := hex.EncodeToString(tmp.Sum(nil))

    rank, err := strconv.ParseBool(admin)
    if err != nil {
        log.Fatalln(err)
    }

    err = db.AddUser(user, hash, rank)
    if err != nil {
        log.Fatalln(err)
    }
    fmt.Println("Success")
}

func rmUser(user string) {
    if err := db.Read(); err != nil {
        log.Fatalln(err)
    }
    if err := db.RemoveUser(user); err != nil {
        log.Fatalln(err)
    }
    fmt.Println("Success")
}

func list() {
    if err := db.Read(); err != nil {
        log.Fatalln(err)
    }
    
    fields := db.QueryAll()
    if len(fields) == 0 {
        fmt.Println("No users found")
        return
    }
    
    for _, field := range fields {
        fmt.Println(field)
    }
}
