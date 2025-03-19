package server

import (
	"backend/internal/types"
	"log"
	"net/http"
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
func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("form.html").Funcs(template.FuncMap{
		"contains": contains,
	}).ParseFiles("./static/form.html"))

	// Получаем данные и ошибки из cookies
	formData, err := getFormDataFromCookies(r)
	if err != nil {
		log.Println(err)
	}
	formErrors, err := getFormErrorsFromCookies(r) // структура ошибок либо nil
	if err != nil {
		log.Println(formErrors)
		log.Println(err)
	}
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
