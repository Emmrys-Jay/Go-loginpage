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
