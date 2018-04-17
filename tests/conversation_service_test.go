package tests

import (
	"../service/conversation_service"
	"../service/message_service"
	"testing"
)

/**************************************************************************
* M O C K  D A T A
**************************************************************************/

var MockConversations map[string]conversation_service.Conversation = map[string]conversation_service.Conversation{
	"ok_1": conversation_service.Conversation{
		UserIds: []string{"someuser1", "someuser2"},
		Messages: []message_service.Message{
			MockMessages["ok_1"],
			MockMessages["ok_2"],
		},
		LastMessageId: MockMessages["ok_2"].Id,
	},
}

/**************************************************************************
* T E S T S
**************************************************************************/

func TestUniqueKey(t *testing.T) {
	// 1. Should be able to get a unique key for a mock conversation
	_conv := MockConversations["ok_1"]
	conv := &_conv
	key := conv.UniqueKey()
	if key != "conversation_someuser1_someuser2" {
		t.Errorf("Invalid conversation key returned, expected %s, got %s", "conversation_someuser1_someuser2", key)
	}
}

func TestAddMessage(t *testing.T) {
	// 1. Adding a message to an existing conversation should work fine
	_conv := MockConversations["ok_1"]
	conv := &_conv
	mId, err := conv.AddMessage(MockMessages["ok_3"])
	if err != nil {
		t.Error(err)
	}
	if mId != 3 {
		t.Errorf("Invalid message id was returned, expected %d, got %d", 3, mId)
	}

	// 2. Adding an empty message should not work
	mId, err = conv.AddMessage(MockMessages["empty_1"])
	if err == nil {
		t.Error("expected the validation to fail on an empty message, but it did not")
	}
}

func TestEditMessage(t *testing.T) {
	// 1. Should be able to edit a message
	_conv := MockConversations["ok_1"]
	conv := &_conv
	err := conv.EditMessage(MockConversations["ok_1"].Messages[1].Id, "edited message", MockConversations["ok_1"].Messages[1].From)
	if err != nil {
		t.Error(err)
	}
	if conv.Messages[1].Content != "edited message" {
		t.Errorf("Editing message failed")
	}
	conv, err = conversation_service.GetConversationByUserIds(MockConversations["ok_1"].UserIds)
	if err != nil {
		t.Error(err)
	}
	if conv.Messages[1].Content != "edited message" {
		t.Errorf("Editing message failed")
	}
}

func TestDeleteMessage(t *testing.T) {
	// 1. Should be able to edit a message
	_conv := MockConversations["ok_1"]
	conv := &_conv
	err := conv.DeleteMessage(MockConversations["ok_1"].Messages[1].Id, MockConversations["ok_1"].Messages[1].From)
	if err != nil {
		t.Error(err)
	}
	if len(MockConversations["ok_1"].Messages) == 1 {
		t.Errorf("Editing message failed, the length didn't change")
	}
	conv, err = conversation_service.GetConversationByUserIds(MockConversations["ok_1"].UserIds)
	if err != nil {
		t.Error(err)
	}
	if len(conv.Messages) != 1 {
		t.Errorf("Editing message failed because unexpected length of messages in the conversation")
	}
}

func TestSave(t *testing.T) {
	// 1. Saving should work fine
	_conv := MockConversations["ok_1"]
	conv := &_conv
	err := conv.Save()
	if err != nil {
		t.Error(err)
	}
}

func TestGetConversationByUserIds(t *testing.T) {
	// 1. We just added two new messages, we should be able to see them
	conv, err := conversation_service.GetConversationByUserIds(MockConversations["ok_1"].UserIds)
	if err != nil {
		t.Error(err)
	}
	if len(conv.Messages) != 2 {
		t.Errorf("Invalid length of messages, expected %d, got %d", 2, len(conv.Messages))
	}
}
