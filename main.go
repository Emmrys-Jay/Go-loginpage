package main

import (
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	template.Must(template.ParseGlob("templates/*"))
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

}

func signup(w http.ResponseWriter, r *http.Request) {

}

func suggest(w http.ResponseWriter, r *http.Request) {

}

func logout(w http.ResponseWriter, r *http.Request) {

}
