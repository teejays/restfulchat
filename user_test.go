package main

import (
	"testing"
)

/**************************************************************************
* I N I T
**************************************************************************/

func initBuddiesMapTest() error {
	exists, err := db.GetStructIfExists(buddiesCollectionName, "buddies_map", &buddiesMap)
	if err != nil {
		return err
	}
	if !exists {
		buddiesMap = make(map[string]map[string]bool)
	}
	return nil
}

/**************************************************************************
* T E S T S
**************************************************************************/

func TestGetUser(t *testing.T) {
	// It should create the User without any issues
	u1, err := GetUser(mockUserId1)
	if err != nil {
		t.Error(err)
	}
	if u1.UserId != mockUserId1 {
		t.Errorf("failed to set the right userId field for the User object")
	}

	// It should trim the whitespace
	u4, err := GetUser(mockUserId4)
	if err != nil {
		t.Error(err)
	}
	if u4.UserId == mockUserId4 {
		t.Errorf("failed to trim the whitespaces")
	}
	if u4.UserId != "testuser4" {
		t.Errorf("wrong user id field in the user object")
	}

	// It should fail
	_, err = GetUser(invalidUserId1)
	if err == nil {
		t.Errorf("failed to give an error on an empty user id")
	}

}

func TestSendMessage(t *testing.T) {
	u1, err := GetUser(mockUserId1)
	if err != nil {
		t.Error(err)
	}
	mId, err := u1.SendMessage(mockUserId2, testMessageContent1)
	if err != nil {
		t.Error(err)
	}
	if mId < 1 {
		t.Errorf("invalid message id returned")
	}
}

func TestGetConversation(t *testing.T) {
	u1, err := GetUser(mockUserId1)
	if err != nil {
		t.Error(err)
	}
	u2, err := GetUser(mockUserId2)
	if err != nil {
		t.Error(err)
	}
	conv, err := u1.GetConversation(u2)
	if err != nil {
		t.Error(err)
	}

	// since we've only sent one message
	if len(conv.UserIds) != 2 {
		t.Errorf("invalid length of userIds in conversation found")
	}
	if len(conv.Messages) != 1 {
		t.Errorf("invalid length of message in conversation found")
	}

}

func TestRegisterBuddy(t *testing.T) {
	u1, err := GetUser(mockUserId1)
	if err != nil {
		t.Error(err)
	}
	u3, err := GetUser(mockUserId3)
	if err != nil {
		t.Error(err)
	}
	err = u1.RegisterBuddy(u3)
	if err != nil {
		t.Error(err)
	}
}

func TestGetBuddies(t *testing.T) {
	u1, err := GetUser(mockUserId1)
	if err != nil {
		t.Error(err)
	}
	buddies, err := u1.GetBuddies()
	if err != nil {
		t.Error(err)
	}
	if len(buddies) != 2 {
		t.Errorf("invalid length of buddies returned, expected 2, got %d", len(buddies))
	}
}
