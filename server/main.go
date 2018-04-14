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
	initUsersPartnersMap()

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
	Message interface{}
}

// GET: Get all the conversations of the user
func getChat(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var userId string = strings.ToLower(p.ByName("userid"))
	if userId == "" {
		log.Fatal("Invalid userid provided")
	}
	partners := getPartnersForUser(userId)
	var data []*Conversation = make([]*Conversation, len(partners))
	for i, partner := range partners {
		data[i] = GetConversation([]string{userId, partner})
	}

	resp := responseStruct{Message: data}
	b, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(b)
}

type PostChatParams struct {
	Content string
	To      string
}

func postChat(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var userId string = strings.ToLower(p.ByName("userid"))
	if userId == "" {
		log.Fatal("Invalid userid provided")
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	var body PostChatParams
	err = json.Unmarshal(b, &body)
	if err != nil {
		log.Fatal(err)
	}
	if body.Content == "" {
		log.Fatal("Empty message")
	}
	body.To = strings.ToLower(body.To)
	if body.To == "" {
		log.Fatal("Invalid recipient")
	}
	if body.To == userId {
		log.Fatal("Cannot send a message to yourself")
	}
	conv := GetConversation([]string{userId, body.To})

	conv.AddMessage(Message{Content: body.Content, From: userId, Timestamp: time.Now()})
	err = conv.Save()
	if err != nil {
		log.Fatalf("[Save Conversation] %v", err)
	}

	addPartnerForUser(userId, body.To)

	resp := responseStruct{Message: "Message sent."}
	b, err = json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(b)
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
