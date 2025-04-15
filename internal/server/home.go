package server

import (
	"backend/internal/types"
	"log"
	"net/http"
	"strings"
	"text/template"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	prevPref := log.Prefix()
	log.SetPrefix(prevPref + " HomeHandler   ")
	defer log.SetPrefix(prevPref)
	tmpl := template.Must(template.New("home.html").Funcs(template.FuncMap{
		"contains": contains,
	}).ParseFiles("./static/home.html"))
	// Получаем данные и ошибки из cookies
	formData, _ := getFormDataFromCookies(r)

	date := strings.Split(formData.Date, "T")
	formData.Date = date[0]

	formErrors, _ := getFormErrorsFromCookies(r) // структура ошибок либо nil
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
