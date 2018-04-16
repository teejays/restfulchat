package conversation_service

import (
	"../message_service"
	"fmt"
	"github.com/teejays/gofiledb"
	"sort"
	"strings"
)

/**************************************************************************
* C O N V E R S A T I O N
**************************************************************************/

// A conversation is made of ordered messages between two Users

type Conversation struct {
	UserIds       []string
	Messages      []message_service.Message
	LastMessageId int
}

var conversationCollectionName string = "conversation"

func GetConversationByUserIds(userIds []string) (*Conversation, error) {
	var c Conversation

	key := uniqueConversationKey(userIds)

	db := gofiledb.GetClient()
	exists, err := db.GetStructIfExists(conversationCollectionName, key, &c)
	if err != nil {
		return nil, err
	}
	if !exists {
		c.UserIds = userIds
	}
	return &c, nil
}

func (c *Conversation) AddMessage(m message_service.Message) (int, error) {
	err := m.Validate()
	if err != nil {
		return -1, err
	}

	m.Id = c.LastMessageId + 1
	c.LastMessageId++

	c.Messages = append(c.Messages, m)
	err = c.Save()
	if err != nil {
		return -1, err
	}
	return m.Id, nil
}

func (c *Conversation) EditMessage(messageId int, content string, from string) error {
	if len(c.Messages) < 1 {
		return fmt.Errorf("No messages found")
	}

	var message *message_service.Message
	for i := 0; i < len(c.Messages); i++ {
		if c.Messages[i].Id == messageId && c.Messages[i].From == from {
			message = &c.Messages[i]
			break
		}
	}

	if message == nil {
		return fmt.Errorf("No message found with the given params")
	}

	err := message.Edit(content)
	if err != nil {
		return err
	}

	err = message.Validate()
	if err != nil {
		return err
	}

	err = c.Save()
	if err != nil {
		return err
	}

	return nil
}

func (c *Conversation) DeleteMessage(messageId int, from string) error {
	if len(c.Messages) < 1 {
		return fmt.Errorf("No messages found")
	}

	var messageExists bool
	var messageIndex int
	for i := 0; i < len(c.Messages); i++ {
		if c.Messages[i].Id == messageId && c.Messages[i].From == from {
			messageExists = true
			messageIndex = i
		}
	}

	if !messageExists {
		return fmt.Errorf("No message found with the given params")
	}

	c.Messages = append(c.Messages[0:messageIndex], c.Messages[messageIndex+1:]...)

	err := c.Save()
	if err != nil {
		return err
	}

	return nil
}

func (c *Conversation) Save() error {
	db := gofiledb.GetClient()
	return db.SetStruct(conversationCollectionName, c.UniqueKey(), c)
}

func (c *Conversation) UniqueKey() string {
	return uniqueConversationKey(c.UserIds)
}

/**************************************************************************
* H E L P E R S
**************************************************************************/

func uniqueConversationKey(userIds []string) string {
	sort.Sort(sort.StringSlice(userIds))
	var key string
	key += "conversation_"
	key += strings.Join(userIds, "_")
	return key
}
