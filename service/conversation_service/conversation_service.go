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

/*
Conversation: A conversation is stored communication between two or more users.
Structure:
-- UserIds: an array of user ids of all the users that are a part of a conversation.
-- Messages: An array of Message between the UserIds, ordered with the oldest up first.
-- LastMessageId (int): Keeps track of the last (also largest) unique message id so the new messages can be given an appropriate id.
*/

/* How are the conversations stored in the DB?
-- Each individual conversation is stored as a separate file
-- Name of the files follow the pattern: "conversation_<userid>_<userid>"
-- ^ The filename is also the "key" that the DB client uses while storing an object
*/

// Define the structure for the Conversation object
type Conversation struct {
	UserIds       []string
	Messages      []message_service.Message
	LastMessageId int
}

// Since we store the conversations in the database, we need to have a collection name it.
var conversationCollectionName string = "conversation"

// Given a list of user ids, load and return the conversation between them
func GetConversationByUserIds(userIds []string) (*Conversation, error) {

	// Initialize an empty conversation variable so we can load the saved conversation file into it
	var c Conversation

	// First we need to get the key used to store the conversation
	key := uniqueConversationKey(userIds)

	// Get the db client and then get the conversation using the key
	db := gofiledb.GetClient()
	exists, err := db.GetStructIfExists(conversationCollectionName, key, &c)
	if err != nil {
		return nil, err
	}
	// If the conversation doesn't exist between the given users
	// Return the empty conversation object after setting the UserIds field
	if !exists {
		c.UserIds = userIds
	}
	return &c, nil
}

// Given a conversation, add a new message to it.
func (c *Conversation) AddMessage(m message_service.Message) (int, error) {
	// Make sure that the message is valid, passes sanity checks
	err := m.Validate()
	if err != nil {
		return -1, err
	}

	// Assign a new message id to the message
	// The new message id is the message id of the last added message + 1
	// We store the message id of the last added message in the LastMessageId field in Conversation
	m.Id = c.LastMessageId + 1
	c.LastMessageId++

	// Append the new message to the conversation messages
	c.Messages = append(c.Messages, m)

	// Save the conversation
	err = c.Save()
	if err != nil {
		return -1, err
	}
	return m.Id, nil
}

// Given a conversation, a sender, and a message id, edit the message to the new content
func (c *Conversation) EditMessage(messageId int, newContent string, from string) error {

	// Before we can edit the message, we need to find it
	// Loop through all the messages to find a message that has the provided message id, and is sent by the provided user id
	var message *message_service.Message
	for i := 0; i < len(c.Messages); i++ {
		if c.Messages[i].Id == messageId && c.Messages[i].From == from {
			// If we find such message, we should store a reference to it, and end the loop
			message = &c.Messages[i]
			break
		}
	}

	// If the message we are looking for is not found in the above loop, return an error
	if message == nil {
		return fmt.Errorf("No message found in the conversation with message id %d", messageId)
	}

	// If the message is found, edit the content of the message
	err := message.Edit(newContent)
	if err != nil {
		return err
	}

	// Save the conversation with the edited message
	err = c.Save()
	if err != nil {
		return err
	}

	return nil
}

// Given a conversation, a sender, and a message id, delete that message from the record
func (c *Conversation) DeleteMessage(messageId int, from string) error {

	// Before we can delete the message, we need to find it
	// Loop through all the messages to find a message that has the provided message id, and is sent by the provided user id
	var messageExists bool
	var messageIndex int
	for i := 0; i < len(c.Messages); i++ {
		if c.Messages[i].Id == messageId && c.Messages[i].From == from {
			// If we find such message, we should store a reference to it, and end the loop
			messageExists = true
			messageIndex = i
			break
		}
	}

	// If the message we're looking for is not found, return an error
	if !messageExists {
		return fmt.Errorf("No message found with the given params")
	}

	// If found, remove it from the message array of the conversation
	c.Messages = append(c.Messages[0:messageIndex], c.Messages[messageIndex+1:]...)

	// Save the conversation
	err := c.Save()
	if err != nil {
		return err
	}

	return nil
}

// Given a conversation object, saves the conversation to the database
func (c *Conversation) Save() error {
	db := gofiledb.GetClient()
	return db.SetStruct(conversationCollectionName, c.UniqueKey(), c)
}

// Given a conversatoin object, returns the unique key that is used to refer to the object while saving and loading in the database
func (c *Conversation) UniqueKey() string {
	return uniqueConversationKey(c.UserIds)
}

/**************************************************************************
* H E L P E R S
**************************************************************************/

// Creates the unique for a conversation between given user ids
// Follows the simple pattern: "conversation_<userid>_<userid>"
func uniqueConversationKey(userIds []string) string {
	sort.Sort(sort.StringSlice(userIds))
	var key string
	key += "conversation_"
	key += strings.Join(userIds, "_")
	return key
}
