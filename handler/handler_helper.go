package handler

import (
	"../service/user_service"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"strings"
)

/**************************************************************************
* H T T P  H E L P E R  F U N C T I O N S
**************************************************************************/

// Authnetication has not been implemented yet.
// We're just passing user ids as a route param, and verifying that it's valid
// In a real application, we should let users bass basic auth tokens and verify
func authenticateRequest(r *http.Request, p httprouter.Params) (*user_service.User, error) {
	uid := strings.ToLower(p.ByName("userid"))
	if uid == "" {
		return nil, fmt.Errorf("Invalid userid provided")
	}
	return user_service.GetUser(uid)
}

// The standard API response struct for any data that our server might return
type ResponseStruct struct {
	IsError bool
	Data    interface{}
}

// Helps a HTTP handler return any data encoded as a json
func writeData(w http.ResponseWriter, data interface{}) {
	resp := ResponseStruct{Data: data}
	b, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

// Helps a HTTP handler return a 400 (Bad Request)
// To do: We should probably have the status code passed as a parameter, so we can handle more StatusCodes
func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	resp := ResponseStruct{Data: fmt.Sprintf("%s", err.Error()), IsError: true}
	b, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

// Helps unmarshal the request body into the provided struct
func parseBody(r *http.Request, v interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	err = json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	return nil
}
