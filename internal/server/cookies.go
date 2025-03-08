package server

import (
	"backend/internal/types"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func clearCookies(w http.ResponseWriter) {

}

func getFormDataFromCookies(r *http.Request) (types.Form, error) {

}

func getFormErrorsFromCookies(r *http.Request) (types.FormErrors, error) {

}

func getSuccessFromCookies(r *http.Request) bool {

}

func setFormDataCookie(w http.ResponseWriter, json_data []byte) {
	log.Println(string(json_data))
	http.SetCookie(w, &http.Cookie{
		Name:     "values",
		Value:    base64.StdEncoding.EncodeToString(json_data),
		Path:     "/process",
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	})
}

func setErrorsCookie(w http.ResponseWriter, formerrors []byte) {
	log.Println(string(formerrors))
	http.SetCookie(w, &http.Cookie{
		Name:     "errors",
		Value:    base64.StdEncoding.EncodeToString(formerrors),
		Path:     "/process",
		Expires:  time.Now().AddDate(1, 0, 0), // 1 year
		HttpOnly: true,
		Secure:   true,
	})
}

func setSucsessCookie(w http.ResponseWriter) {
	data, _ := json.Marshal(1)
	log.Println(string(data))
	http.SetCookie(w, &http.Cookie{
		Name:     "form_success",
		Value:    string(data),
		Path:     "/process",
		Expires:  time.Now().Add(1 * time.Hour), // 1 час
		HttpOnly: true,
	})
}
