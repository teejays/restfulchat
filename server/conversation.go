package main

import (
	"log"
	"sort"
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
	exists, err := db.GetStructIfExists(key, &c, "/conversations/")
	if err != nil {
		log.Fatalf("[GetConversation] %v", err)
	}
	if !exists {
		c.UserIds = userIds
	}
	return &c
}

func (c *Conversation) AddMessage(m Message) error {
	c.Messages = append(c.Messages, m)
	return nil
}

func (c *Conversation) Save() error {
	return db.SetStruct(c.UniqueKey(), c, "/conversations/")
}

func getUniqueKeyForMembers(userIds []string) string {
	sort.Sort(sort.StringSlice(userIds))
	var key string
	key += "conversation_"
	key += strings.Join(userIds, "_")
	return key
}

func (c *Conversation) UniqueKey() string {
	return getUniqueKeyForMembers(c.UserIds)
}
