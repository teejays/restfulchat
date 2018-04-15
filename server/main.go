package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/teejays/gofiledb"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
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
	gofiledb.InitClient("/home/talhajansari/data/restfulchat")
	db = gofiledb.GetClient()
	err := initBuddiesMap()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the server
	router := httprouter.New()
	router.GET("/v1/chat/:userid", getChat)
	router.POST("/v1/chat/:userid", postChat)
	log.Fatal(http.ListenAndServe(":8080", router))
}

/**************************************************************************
* H A N D L E R S
**************************************************************************/
type responseStruct struct {
	IsError bool
	Data    interface{}
}

// GET: Get all the conversations of the user
func getChat(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user, err := authenticateRequest(r, p)
	if err != nil {
		writeError(w, err)
		return
	}

	buddies := user.GetBuddies()
	var data []*Conversation = make([]*Conversation, len(buddies))
	for i, buddy := range buddies {
		data[i], err = user.GetConversation(buddy.UserId)
		if err != nil {
			writeError(w, err)
			return
		}
	}

	writeData(w, data)
}

type PostChatParams struct {
	Message string
	To      string
}

func postChat(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user, err := authenticateRequest(r, p)
	if err != nil {
		writeError(w, err)
		return
	}

	// Get the posted message and process
	var body PostChatParams
	err = parseBody(r, &body)
	if err != nil {
		writeError(w, err)
		return
	}

	recipient := strings.ToLower(body.To)
	err = validateUserId(recipient)
	if err != nil {
		writeError(w, err)
		return
	}

	conv, err := user.GetConversation(recipient)
	if err != nil {
		writeError(w, err)
		return
	}

	messageId, err := conv.AddMessage(Message{Content: body.Message, From: user.UserId, Timestamp: time.Now()})
	if err != nil {
		writeError(w, err)
		return
	}

	err = user.RegisterBuddy(recipient)
	if err != nil {
		writeError(w, err)
		return
	}

	writeData(w, fmt.Sprintf("Message sent. Message Id: %d", messageId))

}

// This is not being implemented yet, so we're just passing user ids as a route param
func authenticateRequest(r *http.Request, p httprouter.Params) (*User, error) {
	uid := strings.ToLower(p.ByName("userid"))
	if uid == "" {
		return nil, fmt.Errorf("Invalid userid provided")
	}
	return GetUser(uid)
}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf("%s", err.Error())))
}

func writeData(w http.ResponseWriter, data interface{}) {
	resp := responseStruct{Data: data}
	b, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func parseBody(r *http.Request, v interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	err = json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	return nil
}

/**************************************************************************
* M O D E L S
**************************************************************************/
// Need to think of a DB design. A major part of this project is logging conversations.
// It makes sense to organize it by conversations
