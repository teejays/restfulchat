package main

import (
	"log"
	"strings"
	"time"
)

// Message is the primary data object
type Message struct {
	Content   string
	Timestamp time.Time
	From      string
}

// A conversation is made of ordered messages between two Users
type Conversation struct {
	UserIds  []string
	Messages []Message
}

func GetConversation(userIds []string) *Conversation {
	key := getUniqueKeyForMembers(userIds)
	var c Conversation
	_, err := db.GetStructIfExists(key, &c, "/restfulchat/conversations")
	if err != nil {
		log.Fatal(err)
	}
	return &c
}

func getUniqueKeyForMembers(userIds []string) string {
	var key string
	key += "conversation_"
	key += strings.Join(userIds, "_")
	return key
}

func (c *Conversation) UniqueKey() string {
	return getUniqueKeyForMembers(c.UserIds)
}

func (c *Conversation) AddMessage(m Message) error {
	c.Messages = append(c.Messages, m)
	return nil
}
