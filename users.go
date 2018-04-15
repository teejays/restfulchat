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
	err := validateUserId(userId)
	if err != nil {
		return nil, err
	}
	userId = processUserId(userId)

	var u User = User{UserId: userId}

	return &u, nil
}

func (u *User) GetConversations() ([]*Conversation, error) {
	buddies, err := u.GetBuddies()
	if err != nil {
		return nil, err
	}

	var data []*Conversation = make([]*Conversation, len(buddies))

	for i, buddy := range buddies {
		data[i], err = u.GetConversation(buddy)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (u *User) GetBuddies() ([]*User, error) {
	_buddies, exists := buddiesMap[u.UserId]
	if !exists {
		return []*User{}, nil
	}
	var buddies []*User
	for bid, v := range _buddies {
		if v {
			buddy, err := GetUser(bid)
			if err != nil {
				return nil, err
			}
			buddies = append(buddies, buddy)
		}
	}
	return buddies, nil
}

func (u *User) GetConversation(buddy *User) (*Conversation, error) {
	if u.UserId == buddy.UserId {
		return nil, fmt.Errorf("Cannot have a conversation with yourself")
	}

	var userIds []string = []string{u.UserId, buddy.UserId}
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

func (u *User) SendMessage(recipientUserId, content string) (int, error) {
	timestamp := time.Now()

	buddy, err := GetUser(recipientUserId)
	if err != nil {
		return -1, err
	}

	conv, err := u.GetConversation(buddy)
	if err != nil {
		return -1, err
	}

	messageId, err := conv.AddMessage(Message{Content: content, From: u.UserId, Timestamp: timestamp})
	if err != nil {
		return messageId, err
	}

	err = u.RegisterBuddy(buddy)
	if err != nil {
		return -1, err
	}

	return messageId, nil

}

func (u *User) RegisterBuddy(buddy *User) error {
	if _, exists := buddiesMap[u.UserId]; !exists {
		buddiesMap[u.UserId] = make(map[string]bool)
	}
	buddiesMap[u.UserId][buddy.UserId] = true

	if _, exists := buddiesMap[buddy.UserId]; !exists {
		buddiesMap[buddy.UserId] = make(map[string]bool)
	}
	buddiesMap[buddy.UserId][u.UserId] = true

	err := db.SetStruct(buddiesCollectionName, "buddies_map", &buddiesMap)
	if err != nil {
		return err
	}
	return nil
}

/**************************************************************************
* B U D D I E S
**************************************************************************/
// BuddiesMap is like a cache, to quickly look up who has an existing conversation with whom
// BuddiesMap is saved in the DB for persistency

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

// Ensures that the user id is valid
// To do: Ensure that there no special characters
func validateUserId(userId string) error {
	if strings.Trim(userId, " ") == "" {
		return fmt.Errorf("User Id validation failed: empty user id")
	}
	return nil
}

// Standardizes the user id before further processing
func processUserId(userId string) string {
	userId = strings.ToLower(userId)
	userId = strings.Trim(userId, " ")
	return userId
}
