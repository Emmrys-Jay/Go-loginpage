package main

import (
	"html/template"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

var tpl *template.Template

type users struct {
	fname    string
	lname    string
	email    string
	uname    string
	password string
}

var dbsessions = make(map[string]string) //cookie value/UUID to user ID
var dbusers = make(map[string]users)     //user ID to users info

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/suggest", suggest)
	http.HandleFunc("/logout", logout)

	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func signup(w http.ResponseWriter, r *http.Request) {
	//check if user is already logged in
	if alreadyloggedin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//tpl.ExecuteTemplate(w, "signup.gohtml", nil)

	//accept form input
	if r.Method == http.MethodPost {
		f := r.FormValue("fname")
		l := r.FormValue("lname")
		e := r.FormValue("email")
		u := r.FormValue("uname")
		p := r.FormValue("psword")

		//check if username already exists
		if _, ok := dbusers[u]; ok {
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}

		//create new session
		uuID := uuid.NewV4()

		http.SetCookie(w, &http.Cookie{
			Name:  "sessions",
			Value: uuID.String(),
		})

		dbsessions[uuID.String()] = u
		dbusers[u] = users{f, l, e, u, p}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func suggest(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "suggest.gohtml", nil)
}

func logout(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "logout.gohtml", nil)
}
