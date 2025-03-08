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
	http.SetCookie(w, &http.Cookie{
		Name:     "values",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "errors",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "form_success",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func getFormDataFromCookies(r *http.Request) (types.Form, error) {
	formcookies, err := r.Cookie("values")
	if err != nil {
		return types.Form{}, err
	}
	val, _ := base64.StdEncoding.DecodeString(formcookies.Value)
	var formvalues types.Form
	if err = json.Unmarshal(val, &formvalues); err == nil {
		return formvalues, nil
	}
	return types.Form{}, err
}

func getFormErrorsFromCookies(r *http.Request) (types.FormErrors, error) {
	errcookies, err := r.Cookie("errors")
	if err != nil {
		return types.FormErrors{}, err
	}
	errors, _ := base64.StdEncoding.DecodeString(errcookies.Value)
	var formerrors types.FormErrors
	if err = json.Unmarshal(errors, &formerrors); err == nil {
		return formerrors, nil
	}
	return types.FormErrors{}, err
}

func getSuccessFromCookies(r *http.Request) bool {
	_, err := r.Cookie("form_success")
	return err == nil
}

func setFormDataCookie(w http.ResponseWriter, json_data []byte) {
	log.Println(string(json_data))
	http.SetCookie(w, &http.Cookie{
		Name:     "values",
		Value:    base64.StdEncoding.EncodeToString(json_data),
		Expires:  time.Now().AddDate(1, 0, 0), // 1 year
		HttpOnly: true,
	})
}

func setErrorsCookie(w http.ResponseWriter, formerrors []byte) {
	log.Println(string(formerrors))
	http.SetCookie(w, &http.Cookie{
		Name:     "errors",
		Value:    base64.StdEncoding.EncodeToString(formerrors),
		Expires:  time.Now().AddDate(1, 0, 0), // 1 year
		HttpOnly: true,
	})
}

func setSuccsessCookie(w http.ResponseWriter) {
	data, _ := json.Marshal(1)
	log.Println(string(data))
	http.SetCookie(w, &http.Cookie{
		Name:     "form_success",
		Value:    string(data),
		Expires:  time.Now().AddDate(1, 0, 0), // 1 year
		HttpOnly: true,
	})
}
