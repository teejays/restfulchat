package message_service

import (
	"fmt"
	"strings"
	"time"
)

/**************************************************************************
* M E S S A G E
**************************************************************************/

// Message is the primary data object

type Message struct {
	Id               int
	Content          string
	TimestampCreated time.Time
	TimestampUpdated time.Time
	From             string
}

func (m *Message) Validate() error {
	if strings.Trim(m.Content, " ") == "" {
		return fmt.Errorf("Message validation failed: empty message")
	}
	return nil
}

func (m *Message) Edit(content string) error {
	oldContent := m.Content
	m.Content = content
	m.TimestampUpdated = time.Now()

	err := m.Validate()
	if err != nil {
		m.Content = oldContent
		return err
	}

	return nil
}
