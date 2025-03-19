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

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})

}

func removeJwtFromCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "key",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func removeUsernameFromCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func removePasswordFromCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "password",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func removeLoginErrorFromCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "login_error",
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

// if OK returns username, nil i; else return "", error
func getUsernameFromCookies(r *http.Request) (string, error) {
	name, err := r.Cookie("username")
	if err != nil {
		return "", err
	}
	return name.Value, nil
}

// if OK returns password, nil i; else return "", error
func getPasswordFromCookies(r *http.Request) (string, error) {
	password, err := r.Cookie("password")
	if err != nil {
		return "", err
	}
	return password.Value, nil
}

func getLoginErrorFromCookies(r *http.Request) (string, error) {
	message, err := r.Cookie("login_error")
	if err != nil {
		return "", err
	}
	return message.Value, nil
}

func setLoginErrorCookie(w http.ResponseWriter, message string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "login_error",
		Value:    message,
		Expires:  time.Now().Add(10 * time.Minute), // 10 minuts
		HttpOnly: true,
	})
}

func setUsernameCookie(w http.ResponseWriter, username string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    username,
		Expires:  time.Now().AddDate(1, 0, 0), // 1 year
		HttpOnly: true,
	})
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

func setSuccessCookie(w http.ResponseWriter) {
	data, _ := json.Marshal(1)
	log.Println(string(data))
	http.SetCookie(w, &http.Cookie{
		Name:     "form_success",
		Value:    string(data),
		Expires:  time.Now().AddDate(1, 0, 0), // 1 year
		HttpOnly: true,
	})
}

func setPasswordCookie(w http.ResponseWriter, password string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "password",
		Value:    password,
		Expires:  time.Now().AddDate(1, 0, 0), // 1 year
		HttpOnly: true,
	})
}
