package main

import (
	"encoding/json"
	"log"
	"os"
)

/**************************************************************************
* I N I T
**************************************************************************/

type ConfigTest struct {
	HttpServerPort int
	GoFiledbRoot   string
}

var configFileTest string = "settings_test.json"

var configTest ConfigTest

func init() {

	file, err := os.Open(configFileTest)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configTest)
	if err != nil {
		log.Fatal(err)
	}
}

func GetConfigTest() *ConfigTest {
	return &configTest
}
