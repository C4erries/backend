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

	var tmpl *template.Template
	_, err := getUsernameFromCookies(r)
	if err != nil {
		tmpl = template.Must(template.New("profileLogin.html").ParseFiles("./static/profileLogin.html"))

		username, err := getUsernameFromCookies(r)
		if err != nil {
			log.Println(" getUsernameFromCookies" + err.Error())
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
		return
	}
	tmpl = template.Must(template.New("profile.html").Funcs(template.FuncMap{
		"contains": contains,
	}).ParseFiles("./static/profile.html"))
	// Получаем данные и ошибки из cookies
	formData, err := getFormDataFromCookies(r)
	if err != nil {
		log.Println(err)
	}
	date := strings.Split(formData.Date, "T")
	formData.Date = date[0]
	log.Println(formData.Fio)

	formErrors, err := getFormErrorsFromCookies(r) // структура ошибок либо nil
	if err != nil {
		log.Println(formErrors)
		log.Println("getFormErrorsFromCookies" + err.Error())
	}
	success := getSuccessFromCookies(r)

	username, _ := getUsernameFromCookies(r)
	password, err := getPasswordFromCookies(r)
	if err == nil {
		removePasswordFromCookies(w)
	}
	// Удаляем cookies после их использования в случае ошибки
	//if !(success) {
	//	clearCookies(w)
	//}

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
