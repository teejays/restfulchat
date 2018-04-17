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

// This is the entry point of the app. We do two main things here:
// -- 1) We initialize the various parts of the app that need to be initialized before the API can start working
// -- 2) Set up a webserver, and forward the GET, POST, PUT and DELETE requests to the respecting handlers

func main() {

	// 1. Initialize the things we need in order to run the application
	// -- Application settings, such as HTTP port, are provided in a settings.json file.
	// -- Let's load that file into our config
	config.InitConfig("./config/settings.json")
	// -- Start the gofiledb database client, so other services in the app can save and load their objects
	gofiledb.InitClient(config.GetConfig().GoFiledbRoot)
	// -- Load an in-memory (from the db) that keeps track of what users converse with what other users
	err := user_service.LoadBuddiesInfoToMemory()
	if err != nil {
		log.Fatal(err)
	}

	// II. Initialize the server
	// -- We have four endpoints, following the RESTful standard.
	// -- Question: Why do we pass the :userid varaible in the URL?
	// -- Answer: We're passing :userid to dummy for authentication.
	// -- In a real app, I would replace it with basic auth tokens

	// -- Define a new router based on httprouter, that can handle our REST API endpoints
	// -- All these requests are handlers by functions defined in the "handler" package
	router := httprouter.New()
	router.GET("/v1/chat/:userid", handler.GetChatHandler)
	router.POST("/v1/chat/:userid", handler.PostChatHandler)
	router.PUT("/v1/chat/:userid", handler.PutChatHandler)
	router.DELETE("/v1/chat/:userid", handler.DeleteChatHandler)

	// -- Start the server, and listen on the port provided in the config
	fmt.Printf("HTTP Server listening on port %d\n", config.GetConfig().HttpServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.GetConfig().HttpServerPort), router))
}
