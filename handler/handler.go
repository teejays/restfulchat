package handler

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/**************************************************************************
* H A N D L E R S
**************************************************************************/

// GET: Get all the conversations of the user
func GetChatHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

type ChatBodyParams struct {
	MessageId int
	Content   string
	To        string
}

// POST: Send a message to a user
func PostChatHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("POST request to /v1/chat")

	// 1. Authenticate (dummy)
	user, err := authenticateRequest(r, p)
	if err != nil {
		writeError(w, err)
		return
	}

	// 2. Parse Request body
	var body ChatBodyParams
	err = parseBody(r, &body)
	if err != nil {
		writeError(w, err)
		return
	}

	// 3. Logic: Send the message
	messageId, err := user.SendMessage(body.To, body.Content)
	if err != nil {
		writeError(w, err)
		return
	}

	// 4. Serve Response
	writeData(w, fmt.Sprintf("Message Id: %d", messageId))

}

// POST: Edit a message
func PutChatHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("PUT request to /v1/chat")

	// 1. Authenticate (dummy)
	user, err := authenticateRequest(r, p)
	if err != nil {
		writeError(w, err)
		return
	}

	// 2. Parse Request body
	var body ChatBodyParams
	err = parseBody(r, &body)
	if err != nil {
		writeError(w, err)
		return
	}

	// 3. Logic: Update the message
	err = user.EditMessage(body.To, body.MessageId, body.Content)
	if err != nil {
		writeError(w, err)
		return
	}

	// 4. Serve Response
	writeData(w, "Message updated")
}

// DELETE: Delete a message to a user
func DeleteChatHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("DELETE request to /v1/chat")

	// 1. Authenticate (dummy)
	user, err := authenticateRequest(r, p)
	if err != nil {
		writeError(w, err)
		return
	}

	// 2. Parse Request body
	var body ChatBodyParams
	err = parseBody(r, &body)
	if err != nil {
		writeError(w, err)
		return
	}

	// 3. Logic: Update the message
	err = user.DeleteMessage(body.To, body.MessageId)
	if err != nil {
		writeError(w, err)
		return
	}

	// 4. Serve Response
	writeData(w, "Message deleted")
}
