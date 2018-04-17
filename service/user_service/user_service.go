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
/*
User is the main 'Subject' of any activity in this chat API.
This implementation is *very* basic.
For simplicity, we only define one field for a User object, the UserId.
Field UserId of a user is defined a string so we can potentially treat it as unique usernames

User: Represents a user.
-- Structure:
-- -- UserId (string)
-- Buddy: a user that another user is interacts with.
*/

// Define what a user object would look like
type User struct {
	UserId string
}

// Given a user id, get a User Object
// This just validates and cleans the user id, and returns a new User object with the provided user id.
func GetUser(userId string) (*User, error) {
	err := validateUserId(userId)
	if err != nil {
		return nil, err
	}
	userId = processUserId(userId)

	var u User = User{UserId: userId}

	return &u, nil
}

// Given a User, get all the conversations that user has been a part of
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

// Given a User, and another user (buddy) get the conversations between them
func (u *User) GetConversation(buddy *User) (*conversation_service.Conversation, error) {
	// It doesn't make sense to get a conversation between two same users
	if u.UserId == buddy.UserId {
		return nil, fmt.Errorf("Cannot have a conversation with yourself")
	}

	// Create a user id slice to pass to the GetConversationByUserIds function
	var userIds []string = []string{u.UserId, buddy.UserId}

	return conversation_service.GetConversationByUserIds(userIds)
}

// Given a User, send a new message to the provided recipient
func (u *User) SendMessage(recipientUserId, content string) (int, error) {
	// Record the timestamp so we know when the message was sent
	timestamp := time.Now()

	// Get the User object representation of the recipient, since most functions like dealing with User objects instead of user ids
	buddy, err := GetUser(recipientUserId)
	if err != nil {
		return -1, err
	}

	// Get the existing conversation between the two users so we can add a new message
	conv, err := u.GetConversation(buddy)
	if err != nil {
		return -1, err
	}

	// Create a new Message object (see message_service.go) with the new content
	var newMessage message_service.Message = message_service.Message{
		Content:          content,
		From:             u.UserId,
		TimestampCreated: timestamp,
		TimestampUpdated: timestamp,
	}

	// Add the newly created message into the conversation
	messageId, err := conv.AddMessage(newMessage)
	if err != nil {
		return messageId, err
	}

	// For quick lookups, we store an in-memory map (called the buddiesInfoMap) of all the users and the users they have conversed with.
	// In case this is the first message between the two users, let's make sure update the buddies map.
	err = u.SaveBuddyInfo(buddy)
	if err != nil {
		return -1, err
	}

	return messageId, nil

}

// Given a User, edits a message that has been sent by that user previously.
func (u *User) EditMessage(recipientUserId string, messageId int, newContent string) error {

	// Get the User object representation of the recipient, since most functions like dealing with User objects instead of user ids
	buddy, err := GetUser(recipientUserId)
	if err != nil {
		return err
	}

	// Get the existing conversation between the two users so we can add a new message
	conv, err := u.GetConversation(buddy)
	if err != nil {
		return err
	}

	// Use the message id to edit the particular message in the conversation
	err = conv.EditMessage(messageId, newContent, u.UserId)
	if err != nil {
		return err
	}

	return nil

}

func (u *User) DeleteMessage(recipientUserId string, messageId int) error {

	// Get the User object representation of the recipient, since most functions like dealing with User objects instead of user ids
	buddy, err := GetUser(recipientUserId)
	if err != nil {
		return err
	}

	// Get the existing conversation between the two users so we can add a new message
	conv, err := u.GetConversation(buddy)
	if err != nil {
		return err
	}

	// Use the message id to delete the particular message in the conversation
	err = conv.DeleteMessage(messageId, u.UserId)
	if err != nil {
		return err
	}

	return nil

}

/**************************************************************************
* B U D D I E S
**************************************************************************/

/* In order to help other functions quickly get a list of all the users that a given user has conversed with,
we maintain a map mapping user ids to all the other user ids they have conversed with.
We store this in-memory but also save a copy in the database.
*/

// BuddiesMap is like a cache, to quickly look up who has an existing conversation with whom
var buddiesInfoMap map[string]map[string]bool
var buddiesCollectionName string = "buddies" // name of the collection when storing in the db

// Given a user, get all the users that it has conversed with.
func (u *User) GetBuddies() ([]*User, error) {

	// If the user doesn't exist in the buddies map, this means it has never talked to anyone
	// Therefore, return an empty array
	_buddies, exists := buddiesInfoMap[u.UserId]
	if !exists {
		return []*User{}, nil
	}

	// We need to return an array of User object but buddies map stores user ids
	// Loop through all the buddies to create a User object for each, and return the array
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

// Given a user, add a new buddy for that user
func (u *User) SaveBuddyInfo(buddy *User) error {
	// Sanity Check: make sure we're not adding a user as it's own buddy because that's weird
	if u.UserId == buddy.UserId {
		fmt.Errorf("Cannot save oneself as it's own buddy")
	}
	// If the user is new in buddies map, we'll have to initialize it's data structure (a go thing)
	if _, exists := buddiesInfoMap[u.UserId]; !exists {
		buddiesInfoMap[u.UserId] = make(map[string]bool)
	}
	// Add the new user as a buddy
	buddiesInfoMap[u.UserId][buddy.UserId] = true

	// Important! If we're adding user B as a buddy for user A
	// We should also add user A as a buddy for user B (since that makes sense)

	// If the user is new in buddies map, we'll have to initialize it's data structure (a go thing)
	if _, exists := buddiesInfoMap[buddy.UserId]; !exists {
		buddiesInfoMap[buddy.UserId] = make(map[string]bool)
	}
	// Add the new user as a buddy
	buddiesInfoMap[buddy.UserId][u.UserId] = true

	// Save the new buddies map into the database so we don't lose it
	db := gofiledb.GetClient()
	err := db.SetStruct(buddiesCollectionName, "buddies_map", &buddiesInfoMap)
	if err != nil {
		return err
	}
	return nil
}

// Upon start of the application, this function loads the buddies map into memory from the db
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

/**************************************************************************
* H E L P E R
**************************************************************************/

// Given a user id, it ensures that the user id is valid
func validateUserId(userId string) error {
	// To do: Ensure that there no special characters

	if strings.Trim(userId, " ") == "" {
		return fmt.Errorf("User Id validation failed: empty user id")
	}
	return nil
}

// Given a user id, this standardizes the id by cleaning it up
func processUserId(userId string) string {
	userId = strings.ToLower(userId)
	userId = strings.Trim(userId, " ")
	return userId
}
