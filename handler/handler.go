package handler

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/**************************************************************************
* H A N D L E R S
**************************************************************************/

// GET: Listens for requests to serve all the conversations log of a user
func GetChatHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("GET request to /v1/chat")

	// 1. Authenticate (dummy) the requester
	user, err := authenticateRequest(r, p)
	if err != nil {
		writeError(w, err)
		return
	}

	// 2. Logic: Fetch all the conversations of the provided user
	data, err := user.GetConversations()
	if err != nil {
		writeError(w, err)
		return
	}

	// 3. Serve Response
	writeData(w, data)
}

// Define a struct that can be used by POST, PUT and DELETE requests to send body
type ChatBodyParams struct {
	MessageId int
	Content   string
	To        string
}

// POST: Listens for requests to send a message to another user
func PostChatHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("POST request to /v1/chat")

	// 1. Authenticate (dummy)
	user, err := authenticateRequest(r, p)
	if err != nil {
		writeError(w, err)
		return
	}

	// 2. Parse body of the request so we know what the message is, and to whom
	var body ChatBodyParams
	err = parseBody(r, &body)
	if err != nil {
		writeError(w, err)
		return
	}

	// 3. Logic: Send the message from the caller to the provided user
	messageId, err := user.SendMessage(body.To, body.Content)
	if err != nil {
		writeError(w, err)
		return
	}

	// 4. Serve Response
	writeData(w, fmt.Sprintf("Message Id: %d", messageId))

}

// PUT: Listens for requests to edit a particular message
func PutChatHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("PUT request to /v1/chat")

	// 1. Authenticate (dummy)
	user, err := authenticateRequest(r, p)
	if err != nil {
		writeError(w, err)
		return
	}

	// 2. Parse the body of the request so we know what message to edit
	var body ChatBodyParams
	err = parseBody(r, &body)
	if err != nil {
		writeError(w, err)
		return
	}

	// 3. Logic: Edit the message mentioned in the request body
	err = user.EditMessage(body.To, body.MessageId, body.Content)
	if err != nil {
		writeError(w, err)
		return
	}

	// 4. Serve Response
	writeData(w, "Message updated")
}

// DELETE: Listens for requests to delete a particular message sent to a given user
func DeleteChatHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("DELETE request to /v1/chat")

	// 1. Authenticate (dummy)
	user, err := authenticateRequest(r, p)
	if err != nil {
		writeError(w, err)
		return
	}

	// 2. Parse the body of the request so we know what message to delete
	var body ChatBodyParams
	err = parseBody(r, &body)
	if err != nil {
		writeError(w, err)
		return
	}

	// 3. Logic: Delete the message mentioned in the request body
	err = user.DeleteMessage(body.To, body.MessageId)
	if err != nil {
		writeError(w, err)
		return
	}

	// 4. Serve Response
	writeData(w, "Message deleted")
}
