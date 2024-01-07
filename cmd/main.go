package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
    "main/db"
    "main/limits"
    "main/handle"
)

func main() {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func () {
        <-c
        log.Println("Saving cache...")
        db.StoreCache()
        os.Exit(0)
    }()

    port := ":8080"
    args := os.Args[1:]
    if len(args) > 0 {
        if _, err := strconv.Atoi(args[0]); err == nil {
            port = ":" + args[0]
        }
    }

    err := db.Initialize() // run after the config has been parsed
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Server running on port " + port)
    mux := http.NewServeMux()
    mux.Handle("/",         limits.RateLimit(handle.IndexHandler))
    mux.Handle("/upload",   limits.RateLimit(handle.UploadHandler))
    mux.Handle("/login",    limits.RateLimit(handle.LoginHandler))
    mux.Handle("/files",    limits.RateLimit(handle.OptionsHandler))

    // serve static files
    mux.Handle("/static/",
        http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

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
