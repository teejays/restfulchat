package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/teejays/gofiledb"
	"log"
	"net/http"
)

// This is RESTful' Service for Chat. It should have a few endpoints:
// GET /chat/
// POST /chat
// PUT /chat
// DELETE /chat
// 1: lets a client register a username
// 2: lets a user send a message to another username
// 3: pushes messages to the desired user when another user sends a message

/**************************************************************************
* I N I T
**************************************************************************/
var db *gofiledb.Client

func main() {
	// Initialize the DB
	gofiledb.InitClient("/home/talhajansari/data")
	db = gofiledb.GetClient()

	// Initialize the server
	router := httprouter.New()
	router.GET("/v1/chat/:user", getChat)
	router.POST("/v1/chat", postChat)
	log.Fatal(http.ListenAndServe(":8080", router))
}

/**************************************************************************
* H A N D L E R S
**************************************************************************/
type responseStruct struct {
	IsError bool
	Message string
}

func getChat(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintf(w, "Requesting chat logs for %s\n", p.ByName("userid"))
}

func postChat(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return
}

/**************************************************************************
* M O D E L S
**************************************************************************/
// Need to think of a DB design. A major part of this project is logging conversations.
// It makes sense to organize it by conversations

/**************************************************************************
* F A C T O R Y
**************************************************************************/
var userConnMap map[string]interface{}

func startUserSession(username string, connection interface{}) error {
	if _, exists := userConnMap[username]; exists {
		return fmt.Errorf("The username '%s' is in use", username)
	}
	userConnMap[username] = connection
	return nil
}

func endUserSession(username string) error {
	delete(userConnMap, username)
	return nil
}
