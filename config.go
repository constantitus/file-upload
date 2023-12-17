package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"golang.org/x/time/rate"
)

const ConfigFile = "config.toml"

var(
    Config = struct{
        StoragePath string
        DatabasePath string
        CSS string
        UuidDefTTL time.Duration
        UuidLongTTL time.Duration
        LoginCooldown time.Duration
        LoginAttempts int
        Rate rate.Limit
        RateBursts int
        RateCooldown time.Duration
        FilesizeMax int64
        // TODO: SaveCache
    }{ // defaults
        "./",
        "./database.db",
        "style.css",
        time.Duration(2),
        time.Duration(24),
        time.Duration(5),
        4,
        2,
        4,
        time.Duration(1),
        int64(50000),
    }
    Style []byte
)

func init() {
    // TODO: filesystem watcher for the live updating
    ParseConfig()
}

// Parse the config toml file and update the Config struct
func ParseConfig() {
    settings := struct {
        Storage_path string
        Database_path string
        CSS string
        UUID_def_ttl any
        UUID_long_ttl any
        Login_cooldown any
        Login_attempts int
        Rate float64
        Rate_bursts int
        Rate_cooldown any
        Filesize_max any
    }{}
    _, err := toml.DecodeFile(ConfigFile, &settings)
    if err != nil {
        log.Println(err)
        return
    }

    if settings.Storage_path == "" {
        Config.StoragePath = "./" 
    } else {
        Config.StoragePath = settings.Storage_path
    }
    Config.DatabasePath = settings.Database_path
    Config.CSS = settings.CSS
    Config.UuidDefTTL = parseTime(settings.UUID_def_ttl)
    Config.UuidLongTTL = parseTime(settings.UUID_long_ttl)
    Config.LoginCooldown = parseTime(settings.Login_cooldown)
    Config.LoginAttempts = settings.Login_attempts
    Config.Rate = rate.Limit(settings.Rate)
    Config.RateBursts = settings.Rate_bursts
    Config.RateCooldown = parseTime(settings.Rate_cooldown)
    Config.FilesizeMax = parseSize(settings.Filesize_max)

    Style = []byte("<style>")
    css, _ := os.ReadFile(Config.CSS)
    Style = append(Style, css...)
    Style = append(Style, []byte("</style>")...)
}

func parseSize(in any) int64 {
    var s string
    switch in.(type) {
        case string:
        s = in.(string)
        case int: return int64(in.(int))
        default: return int64(0)
    }

    suffixes := []string{"gb", "mb", "kb", "b"}
    s = strings.ToLower(s)

    for _, suf := range suffixes {
        if s, found := strings.CutSuffix(s, suf); found {
            if i, err := strconv.Atoi(s); err == nil {
                switch suf {
                    case "gb":
                    return int64(i << 30)
                    case "mb":
                    return int64(i << 20)
                    case "kb":
                    return int64(i << 10)
                    case "b":
                    return int64(i)
                }
            }
        }
    }

    if i, err := strconv.Atoi(s); err != nil {
        return int64(i)
    }
    return int64(0)
}

func parseTime(in any) time.Duration {
    var s string
    switch in.(type) {
        case string:
        s = in.(string)
        case int: return time.Duration(in.(int))
        default: return time.Duration(0)
    }


    suffixes := []string{"s", "m", "h"}
    s = strings.ToLower(s)

    for _, suf := range suffixes {
        if s, found := strings.CutSuffix(s, suf); found {
            if i, err := strconv.Atoi(s); err == nil {
                switch suf {
                    case "h":
                    return time.Duration(i) * time.Hour
                    case "m":
                    return time.Duration(i) * time.Minute
                    case "s":
                    return time.Duration(i) * time.Second
                }
            }
        }
    }
    if i, err := strconv.Atoi(s); err != nil {
        return time.Duration(i)
    }
    return time.Duration(0)
}
