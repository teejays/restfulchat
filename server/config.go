package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	HttpServerPort int
	GoFiledbRoot   string
}

var configFile string = "settings.json"

var config Config

func init() {

	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
}

func GetConfig() *Config {
	return &config
}
