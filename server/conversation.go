package main

import (
	"log"
	"time"
)

// Message is the primary data object
type Message struct {
	Content   string
	Timestamp time.Time
	From      User
}

// A conversation is made of ordered messages between a members
type Conversation struct {
	Members  []User // always sorted by user id
	Messages []Message
}

func GetConversation(members []User) *Conversation {
	key := getUniqueKeyForMembers(members)
	var c Conversation
	_, err := db.GetStructIfExists(key, &c, "conversations")
	if err != nil {
		log.Fatal(err)
	}
	return &c
}

func getUniqueKeyForMembers(members []User) string {
	var key string
	key += "conversation_"
	key += stringsJoin(members, "_")
	return key
}

func (c *Conversation) UniqueKey() string {
	return getUniqueKeyForMembers(c.Members)
}

func (c *Conversation) AddMessage(m Message) error {
	c.Messages = append(c.Messages, m)
	return nil
}
