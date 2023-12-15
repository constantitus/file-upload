package main

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)


func init() {
    RateLimiter.clients = make(map[string]*client)
    logLimiter.clients = make(map[string]*logClient)

    go func() {
        for {
            time.Sleep(time.Minute)

            RateLimiter.mu.Lock()
            for ip, client := range RateLimiter.clients {
                if time.Since(client.lastSeen) > 3*time.Minute {
                    delete(RateLimiter.clients, ip)
                }
            }
            RateLimiter.mu.Unlock()

            logLimiter.mu.Lock()
            for ip, logClient := range logLimiter.clients {
                if time.Since(logClient.lastSeen) > 3*time.Minute {
                    delete(logLimiter.clients, ip)
                }
            }
            logLimiter.mu.Unlock()
        }
    }()
}


type client struct {
    limiter  *rate.Limiter
    lastSeen time.Time
}

var RateLimiter = struct{
    mu sync.Mutex
    clients map[string]*client
}{}
var limiter = rate.NewLimiter(100, 200)
// Global and per-ip rate limits
// Individual IP's can try up to rate times per second with a token bucket
// of size rate_bursts and a cooldown of rate_cooldown
// Global rate is set too 100, bursts to 100 and instead of having a cooldown,
// it returns 429: too many requests
func RateLimit(next func(writer http.ResponseWriter, request *http.Request)) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // global
        if !limiter.Allow() {
            status := http.StatusTooManyRequests
            http.Error(w, "Too many requests", status)
            return
        }

        // per-ip
        ip, _, err := net.SplitHostPort(r.RemoteAddr)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        RateLimiter.mu.Lock()
        if _, found := RateLimiter.clients[ip]; !found {
            RateLimiter.clients[ip] = &client{limiter: rate.NewLimiter(Conf.Rate, Conf.RateBursts)}
        }
        RateLimiter.clients[ip].lastSeen = time.Now()
        if !RateLimiter.clients[ip].limiter.Allow() {
            time.Sleep(Conf.RateCooldown)
        }
        RateLimiter.mu.Unlock()
        next(w, r)
    })
}


type logClient struct {
    lastSeen time.Time
    attempts int
}
var logLimiter = struct{
    mu sync.Mutex
    clients map[string]*logClient
}{}
// Per-ip login limiter. Sets a cooldown of login_cooldown every
// login_attempts -failed- attempts.
//
// *if the last attempt is successful, it still sets the cooldown*
func CheckLimit(ip string) bool {
    logLimiter.mu.Lock()
    client, found := logLimiter.clients[ip]
    if !found {
        logLimiter.clients[ip] = &logClient{lastSeen: time.Now(), attempts: 0}
        logLimiter.mu.Unlock()
        return true
    }
    if time.Since(client.lastSeen) < Conf.LoginCooldown {
        client.attempts += 1
        if client.attempts >= Conf.LoginAttempts {
            logLimiter.mu.Unlock()
            return false
        }
        client.lastSeen = time.Now()
        logLimiter.mu.Unlock()
        return true
    }
    client.attempts = 0
    client.lastSeen = time.Now()
    logLimiter.mu.Unlock()
    return true
}
