package config

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"golang.org/x/time/rate"
)

const configFile = "config.toml"

var (
    StoragePath = "./"
    DatabasePath = "./database.db"
    UuidDefTTL = time.Duration(2)
    UuidLongTTL = time.Duration(24)
    LoginCooldown = time.Duration(5)
    LoginAttempts = 4
    Rate = rate.Limit(2)
    RateBursts = 4
    RateCooldown = time.Duration(1)
    FilesizeMax = int64(50000)
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
        UUID_def_ttl any
        UUID_long_ttl any
        Login_cooldown any
        Login_attempts int
        Rate float64
        Rate_bursts int
        Rate_cooldown any
        Filesize_max any
    }{}
    _, err := toml.DecodeFile(configFile, &settings)
    if err != nil {
        log.Println(err)
        return
    }

    if settings.Storage_path == "" {
        StoragePath = "./" 
    } else {
        StoragePath = settings.Storage_path
    }
    DatabasePath = settings.Database_path
    UuidDefTTL = parseTime(settings.UUID_def_ttl)
    UuidLongTTL = parseTime(settings.UUID_long_ttl)
    LoginCooldown = parseTime(settings.Login_cooldown)
    LoginAttempts = settings.Login_attempts
    Rate = rate.Limit(settings.Rate)
    RateBursts = settings.Rate_bursts
    RateCooldown = parseTime(settings.Rate_cooldown)
    FilesizeMax = sizeAtoi(settings.Filesize_max)
}

func sizeAtoi(in any) int64 {
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
