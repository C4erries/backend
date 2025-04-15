package server

import (
	"backend/internal/types"
	"log"
	"net/http"
	"strings"
	"text/template"
)

func contains(list []int, value int) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

// render и отправка html клиенту
func profileHandler(w http.ResponseWriter, r *http.Request) {
	prevPref := log.Prefix()
	log.SetPrefix(prevPref + "ProfileHandler ")
	defer log.SetPrefix(prevPref)

	if r.Method != "GET" {
		log.Println(r.URL.String() + " Method: " + r.Method + "is Not Allowed here. Allowed Methods: ")
		http.Error(w, "Method: "+r.Method+"is Not Allowed here. Allowed Methods: "+http.MethodGet, http.StatusMethodNotAllowed)
		return
	}

	var tmpl *template.Template
	_, err := getUsernameFromCookies(r)
	if err != nil {
		tmpl = template.Must(template.New("profileLogin.html").ParseFiles("./static/profileLogin.html"))

		username, err := getUsernameFromCookies(r)
		if err != nil && username == "" {
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
		return
	}
	tmpl = template.Must(template.New("profile.html").Funcs(template.FuncMap{
		"contains": contains,
	}).ParseFiles("./static/profile.html"))
	// Получаем данные и ошибки из cookies
	formData, _ := getFormDataFromCookies(r)
	date := strings.Split(formData.Date, "T")
	formData.Date = date[0]

	formErrors, _ := getFormErrorsFromCookies(r) // структура ошибок либо nil

	clearErrorCookies(w)
	success := getSuccessFromCookies(r)

	username, _ := getUsernameFromCookies(r)
	password, err := getPasswordFromCookies(r)
	if err == nil {
		removePasswordFromCookies(w)
	}
	// Удаляем cookies после их использования в случае ошибки
	if !(success) {
		clearCookies(w)
	}

	// Рендерим шаблон с данными
	tmpl.Execute(w, struct {
		Data     types.Form
		Errors   types.FormErrors
		Success  bool
		Username string
		Password string
	}{
		Data:     formData,
		Errors:   formErrors,
		Success:  success,
		Username: username,
		Password: password,
	})
}
func RedirectToProfile(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	} else {
		http.Redirect(w, r, "/profile", http.StatusTemporaryRedirect)
		return
	}
}
