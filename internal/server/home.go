package server

import (
	"log"
	"net/http"
	"text/template"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./static/home.html"))

	username, err := getUsernameFromCookies(r)
	if err != nil {
		log.Println(err)
	} else if username == "" {
		removeUsernameFromCookies(w)
	}

	loginError, _ := getLoginErrorFromCookies(r)
	removeLoginErrorFromCookies(w)

	tmpl.Execute(w, struct {
		Username   string
		LoginError string
	}{
		Username:   username,
		LoginError: loginError,
	})
}
