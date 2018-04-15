package main

import (
	"github.com/teejays/gofiledb"
)

/**************************************************************************
* I N I T
**************************************************************************/

func init() {
	gofiledb.InitClient(GetConfigTest().GoFiledbRoot)
	db = gofiledb.GetClient()

	// We need to clean the test database
	db.FlushAll() // make sure this is not called for the non-test environment

	initBuddiesMapTest()
}

var mockUserId1 string = "testuser1"
var mockUserId2 string = "testuser2"
var mockUserId3 string = "testuser3"
var mockUserId4 string = "     testuser4 "
var invalidUserId1 string = " "

var testMessageContent1 string = "Hello World 1!"
var testMessageContent2 string = "Hello World 2!"
