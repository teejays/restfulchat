package tests

import (
	"../handler"
	"github.com/julienschmidt/httprouter"
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

func TestGetChatHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req := httptest.NewRequest("GET", "/v1/chat/someuser1", nil)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router := httprouter.New()
	router.GET("/book/:id", handler.GetChatHandler)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}
