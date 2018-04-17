package message_service

import (
	"fmt"
	"strings"
	"time"
)

/**************************************************************************
* M E S S A G E
**************************************************************************/

/*
Message: The most basic data unit that makes a conversation.

-- Structure:
-- -- Id (int): unique identifier of a message within a conversation
-- -- Content (string): the content of a message
-- -- TimestampCreated (time): when the message was sent
-- -- TimestampUpdated (time): when the message was last updated
-- -- From (string): contributed the message in a conversation

*/

// Define the structure for a Message
type Message struct {
	Id               int
	Content          string
	TimestampCreated time.Time
	TimestampUpdated time.Time
	From             string
}

// Given a message, perform sanity checks to make sure it's valid
func (m *Message) Validate() error {
	// If the message content is empty, it's invalid
	if strings.Trim(m.Content, " ") == "" {
		return fmt.Errorf("Message validation failed: empty message")
	}
	return nil
}

// Given a message, edit the contents of the message
func (m *Message) Edit(newContent string) error {

	// Save the old content and timestamp in a temp variable as we might need to revert back to it
	oldContent := m.Content
	oldTimestampUpdated := m.TimestampUpdated

	// Update the content of the message to new content, and the timestamp
	m.Content = newContent
	m.TimestampUpdated = time.Now()

	// Make sure that the edited message is still valid
	err := m.Validate()
	if err != nil {
		// If the the edited message is not valid, revert back to old values
		m.Content = oldContent
		m.TimestampUpdated = oldTimestampUpdated
		return err
	}

	return nil
}
