package main

import (
	"fmt"
	"strings"
	"time"
)

/**************************************************************************
* U S E R
**************************************************************************/
type User struct {
	UserId string
}

func GetUser(userId string) (*User, error) {
	userId = strings.ToLower(userId)
	err := validateUserId(userId)
	if err != nil {
		return nil, err
	}
	var u User = User{UserId: userId}
	return &u, nil
}

func (u *User) GetConversations() ([]*Conversation, error) {
	buddies := u.GetBuddies()
	var data []*Conversation = make([]*Conversation, len(buddies))
	var err error
	for i, buddy := range buddies {
		data[i], err = u.GetConversation(buddy.UserId)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (u *User) GetBuddies() []User {
	_buddies, exists := buddiesMap[u.UserId]
	if !exists {
		return []User{}
	}
	var buddies []User
	for uid, v := range _buddies {
		if v {
			buddies = append(buddies, User{UserId: uid})
		}
	}
	return buddies
}

func (u *User) GetConversation(buddyUserId string) (*Conversation, error) {
	err := validateUserId(buddyUserId)
	if err != nil {
		return nil, err
	}
	if u.UserId == buddyUserId {
		return nil, fmt.Errorf("Cannot have a conversation with yourself")
	}
	var userIds []string = []string{u.UserId, buddyUserId}
	key := uniqueConversationKey(userIds)
	var c Conversation
	exists, err := db.GetStructIfExists(conversationCollectionName, key, &c)
	if err != nil {
		return nil, err
	}
	if !exists {
		c.UserIds = userIds
	}
	return &c, nil
}

func (u *User) RegisterBuddy(buddyUserId string) error {
	if _, exists := buddiesMap[u.UserId]; !exists {
		buddiesMap[u.UserId] = make(map[string]bool)
	}
	buddiesMap[u.UserId][buddyUserId] = true

	if _, exists := buddiesMap[buddyUserId]; !exists {
		buddiesMap[buddyUserId] = make(map[string]bool)
	}
	buddiesMap[buddyUserId][u.UserId] = true

	err := db.SetStruct(buddiesCollectionName, "buddies_map", &buddiesMap)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) SendMessage(recipientUserId, content string) (int, error) {
	timestamp := time.Now()

	recipientUserId = processUserId(recipientUserId)
	err := validateUserId(recipientUserId)
	if err != nil {
		return -1, err
	}

	conv, err := u.GetConversation(recipientUserId)
	if err != nil {
		return -1, err
	}

	messageId, err := conv.AddMessage(Message{Content: content, From: u.UserId, Timestamp: timestamp})
	if err != nil {
		fmt.Println("[Add Message Error]")
		return messageId, err
	}

	err = u.RegisterBuddy(recipientUserId)
	if err != nil {
		return -1, err
	}

	return messageId, nil

}

/**************************************************************************
* B U D D I E S
**************************************************************************/
var buddiesMap map[string]map[string]bool
var buddiesCollectionName string = "buddies"

func initBuddiesMap() error {
	exists, err := db.GetStructIfExists(buddiesCollectionName, "buddies_map", &buddiesMap)
	if err != nil {
		return err
	}
	if !exists {
		buddiesMap = make(map[string]map[string]bool)
	}
	return nil
}

// Ensure that the user id is valid
// To do: Ensure that there no special characters
func validateUserId(userId string) error {
	if strings.Trim(userId, " ") == "" {
		return fmt.Errorf("User Id validation failed: empty user id")
	}
	return nil
}

// Standardize the user id before processing
func processUserId(userId string) string {
	userId = strings.ToLower(userId)
	userId = strings.Trim(userId, " ")
	return userId
}
