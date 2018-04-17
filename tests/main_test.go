package tests

import (
	"../config"
	"../service/user_service"
	"github.com/teejays/gofiledb"
	"log"
)

/**************************************************************************
* I N I T
**************************************************************************/

// During testing, this sets up the application with some basic settings
func init() {
	config.InitConfig("settings_test.json")
	// Test DB: For testing, we use a different DB location since we do not want to interfare with the actual database
	gofiledb.InitClient(config.GetConfig().GoFiledbRoot)

	// Since the Test DB might have data from the last test, we should flush the DB
	db := gofiledb.GetClient()
	db.FlushAll()

	// (Just like actual app) Initialize an in-memory map of what users have talked to what other users
	err := user_service.LoadBuddiesInfoToMemory()
	if err != nil {
		log.Fatal(err)
	}
}

/**************************************************************************
* T E S T
**************************************************************************/

// No tests for main.go
