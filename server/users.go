package main

import (
	"fmt"
	"github.com/teejays/gofiledb"
	"log"
)

var usersPartnersMap map[string]map[string]bool

type User struct {
	UserId string
}

func initUsersPartnersMap() {
	fmt.Println("Initializing Users Model")
	loadUsersPartnersMap()
}

func loadUsersPartnersMap() {
	db := gofiledb.GetClient()
	exists, err := db.GetStructIfExists("users_partners_map", &usersPartnersMap, "/users/")
	if err != nil {
		log.Fatalf("[UsersPartners Load] %v", err)
	}
	if !exists {
		fmt.Println("Initializing UsersPartners Empty map")
		usersPartnersMap = make(map[string]map[string]bool)
	}
}

func getPartnersForUser(userId string) []string {
	partnersMap, exists := usersPartnersMap[userId]
	if !exists {
		return []string{}
	}
	var partners []string
	for p, v := range partnersMap {
		if v {
			partners = append(partners, p)
		}
	}
	return partners
}

func addPartnerForUser(userId, partnerId string) {
	if _, exists := usersPartnersMap[userId]; !exists {
		usersPartnersMap[userId] = make(map[string]bool)
	}
	usersPartnersMap[userId][partnerId] = true

	if _, exists := usersPartnersMap[partnerId]; !exists {
		usersPartnersMap[partnerId] = make(map[string]bool)
	}
	usersPartnersMap[partnerId][userId] = true

	err := db.SetStruct("users_partners_map", &usersPartnersMap, "/users/")
	if err != nil {
		log.Fatal(err)
	}
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
