package server

import (
	"backend/internal/types"
	"net/http"
	"text/template"
)

// render и отправка html клиенту
func render(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("form.html"))

	// Получаем данные и ошибки из cookies
	formData, _ := getFormDataFromCookies(r)
	formErrors, _ := getFormErrorsFromCookies(r)
	success := getSuccessFromCookies(r)

	// Удаляем cookies после их использования
	if formErrors != nil || success {
		clearCookies(w)
	}

	// Рендерим шаблон с данными
	tmpl.Execute(w, struct {
		Data types.Form
		Errors
		Success bool
	}{
		Data:    formData,
		Errors:  formErrors,
		Success: success,
	})
}
