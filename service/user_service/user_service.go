package user_service

import (
	"../conversation_service"
	"../message_service"
	"fmt"
	"github.com/teejays/gofiledb"
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

func (u *User) GetConversations() ([]*conversation_service.Conversation, error) {
	buddies, err := u.GetBuddies()
	if err != nil {
		return nil, err
	}

	var data []*conversation_service.Conversation = make([]*conversation_service.Conversation, len(buddies))

	for i, buddy := range buddies {
		data[i], err = u.GetConversation(buddy)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (u *User) GetConversation(buddy *User) (*conversation_service.Conversation, error) {
	if u.UserId == buddy.UserId {
		return nil, fmt.Errorf("Cannot have a conversation with yourself")
	}

	var userIds []string = []string{u.UserId, buddy.UserId}

	return conversation_service.GetConversationByUserIds(userIds)

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

	var newMessage message_service.Message = message_service.Message{
		Content:          content,
		From:             u.UserId,
		TimestampCreated: timestamp,
		TimestampUpdated: timestamp,
	}

	messageId, err := conv.AddMessage(newMessage)
	if err != nil {
		return messageId, err
	}

	err = u.SaveBuddyInfo(buddy)
	if err != nil {
		return -1, err
	}

	return messageId, nil

}

func (u *User) EditMessage(recipientUserId string, messageId int, content string) error {

	buddy, err := GetUser(recipientUserId)
	if err != nil {
		return err
	}

	if u.UserId == buddy.UserId {
		return fmt.Errorf("Tried to update a message that you sent to yourself. Such messages do not exist.")
	}

	conv, err := u.GetConversation(buddy)
	if err != nil {
		return err
	}

	err = conv.EditMessage(messageId, content, u.UserId)
	if err != nil {
		return err
	}

	return nil

}

func (u *User) DeleteMessage(recipientUserId string, messageId int) error {
	buddy, err := GetUser(recipientUserId)
	if err != nil {
		return err
	}

	if u.UserId == buddy.UserId {
		return fmt.Errorf("Tried to delete a message that you sent to yourself. Such messages do not exist.")
	}

	conv, err := u.GetConversation(buddy)
	if err != nil {
		return err
	}

	err = conv.DeleteMessage(messageId, u.UserId)
	if err != nil {
		return err
	}

	return nil

}

/**************************************************************************
* B U D D I E S
**************************************************************************/

func (u *User) GetBuddies() ([]*User, error) {
	_buddies, exists := buddiesInfoMap[u.UserId]
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

func (u *User) SaveBuddyInfo(buddy *User) error {
	if u.UserId == buddy.UserId {
		fmt.Errorf("Cannot save oneself as it's own buddy")
	}
	if _, exists := buddiesInfoMap[u.UserId]; !exists {
		buddiesInfoMap[u.UserId] = make(map[string]bool)
	}
	buddiesInfoMap[u.UserId][buddy.UserId] = true

	if _, exists := buddiesInfoMap[buddy.UserId]; !exists {
		buddiesInfoMap[buddy.UserId] = make(map[string]bool)
	}
	buddiesInfoMap[buddy.UserId][u.UserId] = true

	db := gofiledb.GetClient()
	err := db.SetStruct(buddiesCollectionName, "buddies_map", &buddiesInfoMap)
	if err != nil {
		return err
	}
	return nil
}

// BuddiesMap is like a cache, to quickly look up who has an existing conversation with whom
// BuddiesMap is saved in the DB for persistency

var buddiesInfoMap map[string]map[string]bool
var buddiesCollectionName string = "buddies"

// This loads into memory the map of buddies from the db
func LoadBuddiesInfoToMemory() error {
	db := gofiledb.GetClient()

	exists, err := db.GetStructIfExists(buddiesCollectionName, "buddies_map", &buddiesInfoMap)
	if err != nil {
		return err
	}
	if !exists {
		buddiesInfoMap = make(map[string]map[string]bool)
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
