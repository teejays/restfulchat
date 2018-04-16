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

// The standard API response struct for any data
type responseStruct struct {
	IsError bool
	Data    interface{}
}

// Helps a HTTP handler return any data encoded as json
func writeData(w http.ResponseWriter, data interface{}) {
	resp := responseStruct{Data: data}
	b, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

// Helps a HTTP handler return a 400 (Bad Request)
// To do: We should probably have the status code passed as a parameter
func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	resp := responseStruct{Data: fmt.Sprintf("%s", err.Error()), IsError: true}
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
