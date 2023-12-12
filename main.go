package main

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func main() {
    InitDB() // run after the config has been parsed
    log.Println("Server running...")
    mux := http.NewServeMux()
    mux.Handle("/",        RateLimit(MainHandler))
    mux.Handle("/upload/", RateLimit(UploadHandler))
    mux.Handle("/login/",  RateLimit(LoginHandler))
    server := http.Server{
        Addr:         ":" + strconv.Itoa(Conf.Port),
        Handler:      mux,
        // ReadTimeout:  10000,
        // WriteTimeout: 10000,
    }
    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}


type client struct {
    limiter  *rate.Limiter
    lastSeen time.Time
}

var (
    mu      sync.Mutex
    clients = make(map[string]*client)
)

func init() {
    go func() {
        for {
            time.Sleep(time.Minute)

            mu.Lock()
            for ip, client := range clients {
                if time.Since(client.lastSeen) > 3*time.Minute {
                    delete(clients, ip)
                }
            }
            mu.Unlock()
        }
    }()
}

func RateLimit(next func(writer http.ResponseWriter, request *http.Request)) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip, _, err := net.SplitHostPort(r.RemoteAddr)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        mu.Lock()
        if _, found := clients[ip]; !found {
            clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
        }
        clients[ip].lastSeen = time.Now()
        if !clients[ip].limiter.Allow() {
            time.Sleep(time.Duration(Conf.Rate_limit) * time.Second)
        }
        mu.Unlock()
        next(w, r)
    })
}
