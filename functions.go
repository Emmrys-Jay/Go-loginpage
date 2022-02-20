package main

import (
	"net/http"
)

func alreadyloggedin(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie("sessions")
	if err != nil {
		return false
	}

	un := dbsessions[c.Value]
	_, ok := dbusers[un]
	return ok
}

func getUser(w http.ResponseWriter, r *http.Request) users {
	c, _ := r.Cookie("sessions")

	uID := dbsessions[c.Value]
	user := dbusers[uID]
	return user
}
