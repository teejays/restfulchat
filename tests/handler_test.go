package tests

import (
	"../handler"
	"bytes"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/teejays/gofiledb"
	"net/http"
	"net/http/httptest"
	"testing"
)

/**************************************************************************
* M O C K  D A T A
**************************************************************************/
var mockGetParams map[string]httprouter.Params = map[string]httprouter.Params{
	"ok_1":    httprouter.Params([]httprouter.Param{{Key: "userid", Value: MockUsers["ok_id_1"]}}),
	"empty_1": httprouter.Params{},
}

/**************************************************************************
* T E S T S
**************************************************************************/

func TestPostChatHandler(t *testing.T) {
	// This is our first handler test, so let's just refresh the database
	db := gofiledb.GetClient()
	db.FlushAll()

	// Create a message to pass to our request
	_body := handler.ChatBodyParams{
		Content: "Hello someuser2 (original)",
		To:      "someuser2",
	}

	body, err := json.Marshal(_body)
	if err != nil {
		t.Error(err)
	}

	// Create a request to pass to our handler.
	req := httptest.NewRequest("POST", "/v1/chat/someuser1", bytes.NewBuffer(body))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router := httprouter.New()
	router.POST("/v1/chat/:userid", handler.PostChatHandler)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Unmarshal the response so we can check it whether it's what we expected
	var resp handler.ResponseStruct
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Error(err)
	}

	// Check whether the response sent IsError set to true
	// This usually means a 4xx code, so should caught in the status code check
	if resp.IsError {
		t.Errorf("Post /v1/chat endpoint sent a response with IsError set to true.")
	}

	// Check if the response data field is what we expect
	expectedData := "Message Id: 1"
	if resp.Data.(string) != expectedData {
		t.Errorf("Post /v1/chat endpoint sent an unexpected response. Expected %s, got %s", expectedData, resp.Data)
	}

}

func TestPutChatHandler(t *testing.T) {
	// Create a message to pass to our request
	_body := handler.ChatBodyParams{
		MessageId: 1,
		Content:   "Hello someuser2 (edited)",
		To:        "someuser2",
	}

	body, err := json.Marshal(_body)
	if err != nil {
		t.Error(err)
	}

	// Create a request to pass to our handler.
	req := httptest.NewRequest("PUT", "/v1/chat/someuser1", bytes.NewBuffer(body))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router := httprouter.New()
	router.PUT("/v1/chat/:userid", handler.PutChatHandler)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Unmarshal the response so we can check it whether it's what we expected
	var resp handler.ResponseStruct
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Error(err)
	}

	// Check whether the response sent IsError set to true
	// This usually means a 4xx code, so should caught in the status code check
	if resp.IsError {
		t.Errorf("Put /v1/chat endpoint sent a response with IsError set to true.")
	}

	// Check if the response data field is what we expect
	expectedData := "Message updated"
	if resp.Data.(string) != expectedData {
		t.Errorf("Put /v1/chat endpoint sent an unexpected response. Expected %s, got %s", expectedData, resp.Data)
	}

}

func TestGetChatHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req := httptest.NewRequest("GET", "/v1/chat/someuser1", nil)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router := httprouter.New()
	router.GET("/v1/chat/:userid", handler.GetChatHandler)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Unmarshal the response so we can check it whether it's what we expected
	var resp handler.ResponseStruct
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Error(err)
	}

	// Check whether the response sent IsError set to true
	// This usually means a 4xx code, so should caught in the status code check
	if resp.IsError {
		t.Errorf("Get /v1/chat endpoint sent a response with IsError set to true.")
	}

	// Check if the response data field is what we expect
	// The data field should be an array (of conversations)
	convs, ok := resp.Data.([]interface{})
	if !ok {
		t.Errorf("Expected the data field in response to be an array, but couldn't assert it as []interface{}")
	}
	// Since we're posted one message in the TestPostHandler method, we should expect:
	// there shoudl be just one conversation
	if len(convs) != 1 {
		t.Errorf("Expected number of conversations to be returned to be %d, but got %d", 1, len(convs))
	}

	// Assert the only conversation as a map[string]interface{} and so forth so we can take a look inside
	conv, ok := convs[0].(map[string]interface{})
	if !ok {
		t.Errorf("Couldn't assert the first element of conversations as a map[string]interface{}")
	}

	// There should be one messages in the conversation
	messages, ok := conv["Messages"].([]interface{})
	if !ok {
		t.Errorf("Couldn't assert the the Message field in the first element of conversations as an []interface{}")
	}
	if len(messages) != 1 {
		t.Errorf("Expected number of messages in the conversation to be %d, got %d", 1, len(messages))
	}

	// There should be two users in UserIds field of the conversation
	userIds, ok := conv["UserIds"].([]interface{})
	if !ok {
		t.Errorf("Couldn't assert the the UserIds field in the first element of conversations as an []interface{}")
	}
	if len(userIds) != 2 {
		t.Errorf("Expected number of users in the conversation to be %d, got %d", 2, len(userIds))
	}

}

func TestDeleteChatHandler(t *testing.T) {
	// Create a message to pass to our request
	_body := handler.ChatBodyParams{
		MessageId: 1,
		To:        "someuser2",
	}

	body, err := json.Marshal(_body)
	if err != nil {
		t.Error(err)
	}

	// Create a request to pass to our handler.
	req := httptest.NewRequest("DELETE", "/v1/chat/someuser1", bytes.NewBuffer(body))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router := httprouter.New()
	router.DELETE("/v1/chat/:userid", handler.DeleteChatHandler)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Unmarshal the response so we can check it whether it's what we expected
	var resp handler.ResponseStruct
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Error(err)
	}

	// Check whether the response sent IsError set to true
	// This usually means a 4xx code, so should caught in the status code check
	if resp.IsError {
		t.Errorf("Put /v1/chat endpoint sent a response with IsError set to true.")
	}

	// Check if the response data field is what we expect
	expectedData := "Message deleted"
	if resp.Data.(string) != expectedData {
		t.Errorf("Delete /v1/chat endpoint sent an unexpected response. Expected %s, got %s", expectedData, resp.Data)
	}

}
