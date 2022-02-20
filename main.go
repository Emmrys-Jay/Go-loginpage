package main

import (
	"html/template"
	"io"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

var tpl *template.Template

type users struct {
	Fname    string
	Lname    string
	Email    string
	Uname    string
	Password string
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
	if !alreadyloggedin(w, r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user := getUser(w, r)

	tpl.ExecuteTemplate(w, "index.gohtml", user)
}

func login(w http.ResponseWriter, r *http.Request) {
	if alreadyloggedin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	if r.Method == http.MethodPost {
		u := r.FormValue("uname")
		p := r.FormValue("psword")

		if user, ok := dbusers[u]; ok {
			if p == user.Password {

				//create new session
				uuID := uuid.NewV4()
				http.SetCookie(w, &http.Cookie{
					Name:  "sessions",
					Value: uuID.String(),
				})
				dbsessions[uuID.String()] = u

				//redirect to dashboard
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			} else {
				w.Header().Set("Content-Type", "text/html; charset=UTF-8")
				io.WriteString(w, `<h5>Invalid Username/Password</h5><a href="/login"> Sign In Again </a><br> <a href="/signup"> Sign Up </a>`)
				return
			}
		} else {
			w.Header().Set("Content-Type", "text/html; charset=UTF-8")
			io.WriteString(w, `<h5>Invalid Username/Password</h5><a href="/login"> Sign In Again </a><br> <a href="/signup"> Sign Up </a>`)
			return
		}
	}
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func signup(w http.ResponseWriter, r *http.Request) {
	//check if user is already logged in
	if alreadyloggedin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

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
		return
	}

	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func suggest(w http.ResponseWriter, r *http.Request) {
	if !alreadyloggedin(w, r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	user := getUser(w, r)

	tpl.ExecuteTemplate(w, "suggest.gohtml", user)
}

func logout(w http.ResponseWriter, r *http.Request) {
	if !alreadyloggedin(w, r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	//get uuid value and delete it's field
	c, _ := r.Cookie("sessions")
	delete(dbsessions, c.Value)
	c.MaxAge = -1 //delete cookie
	tpl.ExecuteTemplate(w, "logout.gohtml", nil)
}
