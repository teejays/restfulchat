package main

import (
	"github.com/teejays/gofiledb"
	"log"
)

var userConversationMap map[string][]string

type User struct {
	UserId string
}

func init() {
	loadUserConversationMap()
}

func loadUserConversationMap() {
	db := gofiledb.GetClient()
	exists, err := db.GetStructIfExists("user_conversation_map", &userConversationMap, "/restfulchat")
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		userConversationMap = make(map[string][]string)
	}
}

func getConversationPartnersForUser(userId string) []string {
	convs, exists := userConversationMap[userId]
	if !exists {
		return []string{}
	}
	return convs
}

// Takes an array of Users and joins them
func stringsJoin(users []User, sep string) string {
	var str string
	for i, u := range users {
		str += u.UserId
		if i < len(users)-1 {
			str += sep
		}
	}
	return str
}
