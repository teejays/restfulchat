package tests

import (
	"../service/user_service"
	"testing"
)

/**************************************************************************
* M O C K  D A T A
**************************************************************************/

// Mock Data
var MockUsers map[string]string = map[string]string{
	"ok_id_1":         "testuser1",
	"ok_id_2":         "testuser2",
	"ok_id_3":         "testuser3",
	"empty_id_1":      "  ",
	"whitespace_id_1": "   whitespaceuser1  ",
}

var MockContent map[string]string = map[string]string{
	"ok_1":    "Hello World 1!",
	"ok_2":    "Hello World 2!",
	"ok_3":    "Hello World 3!",
	"empty_1": "  ",
}

/**************************************************************************
* T E S T S
**************************************************************************/

func TestGetUser(t *testing.T) {
	// 1. Good UserId should work
	u, err := user_service.GetUser(MockUsers["ok_id_1"])
	if err != nil {
		t.Error(err)
	}
	if u.UserId != MockUsers["ok_id_1"] {
		t.Errorf("Failed to set the right userId field for the User object")
	}

	// 2. UserId with whitespaces on the edges should be trimmed
	u, err = user_service.GetUser(MockUsers["whitespace_id_1"])
	if err != nil {
		t.Error(err)
	}
	if u.UserId == MockUsers["whitespace_id_1"] {
		t.Errorf("Failed to trim the whitespaces")
	}
	if u.UserId != "whitespaceuser1" {
		t.Errorf("Wrong user id field in the user object")
	}

	// 3. Empty UserId should fail
	_, err = user_service.GetUser(MockUsers["empty_id_1"])
	if err == nil {
		t.Errorf("Failed to give an error on an empty user id")
	}

}

func TestSendMessage(t *testing.T) {
	// 1. Good content for good user and good recipient should work
	u, err := user_service.GetUser(MockUsers["ok_id_1"])
	if err != nil {
		t.Error(err)
	}
	mId, err := u.SendMessage(MockUsers["ok_id_2"], MockContent["ok_1"])
	if err != nil {
		t.Error(err)
	}
	if mId < 1 {
		t.Errorf("Invalid message id returned")
	}

	buddy, err := user_service.GetUser(MockUsers["ok_id_2"])
	if err != nil {
		t.Error(err)
	}
	conv, err := u.GetConversation(buddy)
	if err != nil {
		t.Error(err)
	}
	if len(conv.Messages) < 1 {
		t.Errorf("Zero messages in the conversation even after message was sent")
	}
	if conv.Messages[len(conv.Messages)-1].Content != MockContent["ok_1"] {
		t.Errorf("Message that was just sent was not found as the most recent in the conversation")
	}

	// 2. Sending message to one self shouldn't work
	_, err = u.SendMessage(MockUsers["ok_id_1"], MockContent["ok_1"])
	if err == nil {
		t.Errorf("SendMessage() allowed the sending of message from the same user and recipient")
	}
}

func TestGetConversation(t *testing.T) {
	// 1. Since we just tested sending of messages, we should get a conversation
	// -- This conversation should have two users and one message
	u, err := user_service.GetUser(MockUsers["ok_id_1"])
	if err != nil {
		t.Error(err)
	}
	buddy, err := user_service.GetUser(MockUsers["ok_id_2"])
	if err != nil {
		t.Error(err)
	}
	conv, err := u.GetConversation(buddy)
	if err != nil {
		t.Error(err)
	}

	if len(conv.UserIds) != 2 {
		t.Errorf("Invalid length of userIds in conversation found")
	}
	if len(conv.Messages) != 1 {
		t.Errorf("Invalid length of message in conversation found")
	}
	if conv.LastMessageId != 1 {
		t.Errorf("Invalid LastMessageId in conversation found")
	}

	// 2. Trying to get a conversation of two new users should return empty conversation
	buddy, err = user_service.GetUser(MockUsers["ok_id_3"])
	if err != nil {
		t.Error(err)
	}
	conv, err = u.GetConversation(buddy)
	if err != nil {
		t.Error(err)
	}
	if len(conv.UserIds) != 2 {
		t.Errorf("Invalid length of userIds in conversation found for a new conversation")
	}
	if len(conv.Messages) != 0 {
		t.Errorf("Invalid length of message in conversation found for a new conversation")
	}
	if conv.LastMessageId != 0 {
		t.Errorf("Invalid LastMessageId in conversation found for a new conversation")
	}

}

func TestGetBuddies(t *testing.T) {
	// 1. Since we just tested a conversation above, those two users should be saved as buddies
	u, err := user_service.GetUser(MockUsers["ok_id_2"])
	if err != nil {
		t.Error(err)
	}
	buddies, err := u.GetBuddies()
	if err != nil {
		t.Error(err)
	}
	if len(buddies) != 1 {
		t.Errorf("Invalid length of buddies returned, expected %d, got %d", 1, len(buddies))
	}
	if buddies[0].UserId != MockUsers["ok_id_1"] {
		t.Errorf("Unexpected user id for buddy found, expected %s, got %s", MockUsers["ok_id_1"], buddies[0])
	}
}
