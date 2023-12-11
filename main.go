package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
    fmt.Println("Server running...")
    http.HandleFunc("/", MainHandler)
    http.HandleFunc("/upload/", UploadHandler)
    http.HandleFunc("/login/", LoginHandler)
    server := http.Server{
        Addr:         ":" + strconv.Itoa(Conf.Port),
        Handler:      nil,
        // ReadTimeout:  10000,
        // WriteTimeout: 10000,
    }
    log.Fatal(server.ListenAndServe())
}
