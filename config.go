package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

const ConfigFile = "config.toml"

type Config struct {
    Port int
    Database_path string
    Storage_path string
    Default_ttl int
    Rememberme_ttl int
}
var Conf Config

func init() {
    // defaults
    Conf = Config{
        8080,
        "./database.db",
        "./",
        2,
        24,
    }
    _, err := toml.DecodeFile(ConfigFile, &Conf)
    if err != nil {
        log.Println(err)
    }
    if Conf.Storage_path == "" { Conf.Storage_path = "./" }
}
