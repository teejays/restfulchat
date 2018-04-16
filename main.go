package main

import (
	"./config"
	"./handler"
	"./service/user_service"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/teejays/gofiledb"
	"log"
	"net/http"
)

/**************************************************************************
* E N T R Y  P O I N T
**************************************************************************/

func main() {

	// I. Initialize the things we need in order to run the application
	// -- Initialize the application config
	config.InitConfig("./config/settings.json")
	// -- Initialize the database client
	fmt.Println("initializing the client")
	gofiledb.InitClient(config.GetConfig().GoFiledbRoot)
	// -- Initialize an in-memory map of what users have talked to what other users
	err := user_service.LoadBuddiesInfoToMemory()
	if err != nil {
		log.Fatal(err)
	}

	// II. Initialize the server
	router := httprouter.New()
	router.GET("/v1/chat/:userid", handler.GetChatHandler)
	router.POST("/v1/chat/:userid", handler.PostChatHandler)
	router.PUT("/v1/chat/:userid", handler.PutChatHandler)
	router.DELETE("/v1/chat/:userid", handler.DeleteChatHandler)

	fmt.Printf("HTTP Server listening on port %d\n", config.GetConfig().HttpServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.GetConfig().HttpServerPort), router))
}
