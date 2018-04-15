package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

var conversationCollectionName string = "conversation"

// Message is the primary data object
type Message struct {
	Id        int
	Content   string
	Timestamp time.Time
	From      string
}

// A conversation is made of ordered messages between two Users
type Conversation struct {
	UserIds  []string
	Messages []Message
}

func (c *Conversation) AddMessage(m Message) (int, error) {
	err := m.Validate()
	if err != nil {
		return -1, err
	}
	m.Id = len(c.Messages) + 1
	c.Messages = append(c.Messages, m)
	err = c.Save()
	if err != nil {
		return -1, err
	}
	return m.Id, nil
}

func (c *Conversation) Save() error {
	return db.SetStruct(conversationCollectionName, c.UniqueKey(), c)
}

func (m *Message) Validate() error {
	if strings.Trim(m.Content, " ") == "" {
		return fmt.Errorf("Message validation failed: empty message")
	}
	return nil
}

func (c *Conversation) UniqueKey() string {
	return uniqueConversationKey(c.UserIds)
}

func uniqueConversationKey(userIds []string) string {
	sort.Sort(sort.StringSlice(userIds))
	var key string
	key += "conversation_"
	key += strings.Join(userIds, "_")
	return key
}
