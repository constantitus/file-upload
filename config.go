package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

const ConfigFile = "config.toml"

var Conf = struct{
    Port int
    // CSS string
    Storage_path string
    Database_path string
    UUID_def_ttl int
    UUID_long_ttl int
    Login_ttl int
    Rate_limit int
}{
        8080,
        // "style.css",
        "./",
        "./database.db",
        2,
        24,
        5,
        1,
    }

func init() {
    _, err := toml.DecodeFile(ConfigFile, &Conf)
    if err != nil {
        log.Println(err)
    }
    if Conf.Storage_path == "" { Conf.Storage_path = "./" }
}
