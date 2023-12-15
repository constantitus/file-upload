package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
    port := ":8080"
    args := os.Args[1:]
    if len(args) > 0 {
        if _, err := strconv.Atoi(args[0]); err == nil {
            port = ":" + args[0]
        }
    }
    InitDB() // run after the config has been parsed
    log.Println("Server running on port " + port)
    mux := http.NewServeMux()
    mux.Handle("/",         RateLimit(MainHandler))
    mux.Handle("/upload/",  RateLimit(UploadHandler))
    mux.Handle("/login/",   RateLimit(LoginHandler))
    server := http.Server{
        Addr:         port,
        Handler:      mux,
        // ReadTimeout:  10000,
        // WriteTimeout: 10000,
    }
    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}
