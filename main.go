package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/teejays/gofiledb"
	"log"
	"net/http"
)

// This is RESTful' Service for Chat. It should have a few endpoints:
// GET /chat/ [DONE]
// POST /chat [DONE]
// PUT /chat [To do]
// DELETE /chat [To do]

/**************************************************************************
* I N I T
**************************************************************************/
var db *gofiledb.Client

func main() {
	// I. Initialize the things we need in order to run the application
	// a) Initialize the database client
	gofiledb.InitClient(GetConfig().GoFiledbRoot)
	db = gofiledb.GetClient()
	// b) Initialize an in-memory map of what users have talked to what other users
	err := initBuddiesMap()
	if err != nil {
		log.Fatal(err)
	}

	// II. Initialize the server
	router := httprouter.New()
	router.GET("/v1/chat/:userid", getChatHandler)
	router.POST("/v1/chat/:userid", postChatHandler)
	fmt.Printf("HTTP Server listening on port %d\n", GetConfig().HttpServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", GetConfig().HttpServerPort), router))
}

/**************************************************************************
* H A N D L E R S
**************************************************************************/

// GET: Get all the conversations of the user
func getChatHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("GET request to /v1/chat")

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
	fmt.Println("POST request to /v1/chat")

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
