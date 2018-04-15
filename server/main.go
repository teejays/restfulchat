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
	gofiledb.InitClient(GetConfig().GoFiledbRoot)
	db = gofiledb.GetClient()
	err := initBuddiesMap()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the server
	router := httprouter.New()
	router.GET("/v1/chat/:userid", getChatHandler)
	router.POST("/v1/chat/:userid", postChatHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", GetConfig().HttpServerPort), router))
}

/**************************************************************************
* H A N D L E R S
**************************************************************************/

// GET: Get all the conversations of the user
func getChatHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// 1. Authenticate (dummy) the requester
	user, err := authenticateRequest(r, p)
	if err != nil {
		writeError(w, err)
		return
	}

	// 2. Do the logic
	data, err := user.GetConversations()
	if err != nil {
		writeError(w, err)
		return
	}

	// 3. Serve Response
	writeData(w, data)
}

type PostChatParams struct {
	Message string
	To      string
}

// POST: Send a message to a user
func postChatHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// 1. Authenticate (dummy)
	user, err := authenticateRequest(r, p)
	if err != nil {
		writeError(w, err)
		return
	}

	// 2. Parse Request body
	var body PostChatParams
	err = parseBody(r, &body)
	if err != nil {
		writeError(w, err)
		return
	}

	// 3. Logic: Send the message
	messageId, err := user.SendMessage(body.To, body.Message)
	if err != nil {
		writeError(w, err)
		return
	}

	// 4. Serve Response
	writeData(w, fmt.Sprintf("Message Id: %d", messageId))

}

/**************************************************************************
* M O D E L S
**************************************************************************/
// Need to think of a DB design. A major part of this project is logging conversations.
// It makes sense to organize it by conversations
