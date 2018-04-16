package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	HttpServerPort int
	GoFiledbRoot   string
}

var config Config

func InitConfig(configFilePath string) {

	file, err := os.Open(configFilePath)
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
