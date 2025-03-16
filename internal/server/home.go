package server

import (
	"net/http"
	"text/template"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("home.html")

	username, _ := getUserNameFromCookies(r)
	success := getSuccessFromCookies(r)
	tmpl.Execute(w, struct {
		Username string
		Success  bool
	}{
		Username: username,
		Success:  success,
	})
}
