package tests

import (
	"../service/message_service"
	"testing"
	"time"
)

/**************************************************************************
* M O C K  D A T A
**************************************************************************/

var MockMessages map[string]message_service.Message = map[string]message_service.Message{
	"ok_1": message_service.Message{
		Id:               1,
		Content:          "Hello world 1",
		TimestampCreated: time.Now(),
		TimestampUpdated: time.Now(),
		From:             "someuser1",
	},
	"ok_2": message_service.Message{
		Id:               2,
		Content:          "Hello world 2",
		TimestampCreated: time.Now(),
		TimestampUpdated: time.Now(),
		From:             "someuser2",
	},
	"ok_3": message_service.Message{
		Id:               3,
		Content:          "Hello world 3",
		TimestampCreated: time.Now(),
		TimestampUpdated: time.Now(),
		From:             "someuser1",
	},
	"empty_1": message_service.Message{
		Id:               4,
		Content:          "    ",
		TimestampCreated: time.Now(),
		TimestampUpdated: time.Now(),
		From:             "someuser2",
	},
}

/**************************************************************************
* T E S T S
**************************************************************************/

func TestValidate(t *testing.T) {
	// 1. Ok message should be fine
	_m := MockMessages["ok_1"]
	m := &_m
	err := m.Validate()
	if err != nil {
		t.Error("A valid message was not validated")
	}
	// 2. Empty message should not be validated
	_m = MockMessages["empty_1"]
	err = m.Validate()
	if err == nil {
		t.Errorf("An invalid message was validated")
	}

}
