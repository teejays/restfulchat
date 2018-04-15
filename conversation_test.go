package main

import (
	"testing"
	"time"
)

/**************************************************************************
* I N I T
**************************************************************************/

/**************************************************************************
* T E S T S
**************************************************************************/

var mockMessage1 Message = Message{
	Id:        1,
	Content:   "Hello worlds 1",
	Timestamp: time.Now(),
	From:      mockUserId2,
}

var mockMessage2 Message = Message{
	Id:        2,
	Content:   "Hello worlds 2",
	Timestamp: time.Now(),
	From:      mockUserId3,
}

var newMessage1 Message = Message{
	Id:        0,
	Content:   "Hello worlds 3",
	Timestamp: time.Now(),
	From:      mockUserId3,
}

var invalidMessage1 Message = Message{
	Id:        1,
	Content:   "   ",
	Timestamp: time.Now(),
	From:      mockUserId2,
}

var mockConversation1 Conversation = Conversation{
	UserIds:  []string{mockUserId2, mockUserId3},
	Messages: []Message{mockMessage1, mockMessage2},
}

func TestUniqueConversationKey(t *testing.T) {

	key := uniqueConversationKey([]string{mockUserId3, mockUserId2})
	if key != "conversation_testuser2_testuser3" {
		t.Errorf("invalid conversation key returned, expected %s, got %s", "conversation_testuser1_testuser2", key)
	}
}

func TestUniqueKey(t *testing.T) {
	key := mockConversation1.UniqueKey()
	if key != "conversation_testuser2_testuser3" {
		t.Errorf("invalid conversation key returned, expected %s, got %s", "conversation_testuser2_testuser3", key)
	}
}

func TestValidate(t *testing.T) {
	err := invalidMessage1.Validate()
	if err == nil {
		t.Errorf("an invalid message was validated, which is weird")
	}
}

func TestAddMessage(t *testing.T) {
	mId, err := mockConversation1.AddMessage(newMessage1)
	if err != nil {
		t.Error(err)
	}
	if mId != 3 {
		t.Errorf("an invalid message id was returned, expected 3, got %d", mId)
	}
}

func TestSave(t *testing.T) {
	err := mockConversation1.Save()
	if err != nil {
		t.Error(err)
	}
}
