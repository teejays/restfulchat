package main

type User string

func stringsJoin(users []User, sep string) string {
	var str string
	for i, u := range users {
		str += string(u)
		if i < len(users)-1 {
			str += sep
		}
	}
	return str
}
