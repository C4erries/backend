package server

import (
	"backend/internal/database"
	"backend/internal/types"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type JsonForm struct {
	types.Form
	Familiar string `json:"Familiar"`
}

func castToJsonForm(values url.Values) *JsonForm {
	var favlangs []int
	for _, value := range values["Favlangs"] {
		num, _ := strconv.Atoi(value)
		favlangs = append(favlangs, num)
	}
	return &JsonForm{
		types.Form{
			Fio:      values["Fio"][0],
			Tel:      values["Tel"][0],
			Email:    values["Email"][0],
			Date:     values["Date"][0],
			Gender:   values["Gender"][0],
			Favlangs: favlangs,
			Bio:      values["Bio"][0],
		},
		values["Familiar"][0]}
}

func processRequestParser(r *http.Request) (unvalidatedForm *JsonForm, err error) {
	// Определяем Content-Type
	contentType := r.Header.Get("Content-Type")

	// Парсим JSON
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&unvalidatedForm); err != nil {
			return nil, errors.New("invalid JSON format: " + err.Error())
		}
		return unvalidatedForm, nil
	}

	// Парсим форму
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			return nil, errors.New("failed to parse form")
		}
		return castToJsonForm(r.Form), nil
	}

	return nil, errors.New("unsupported content type")
}

func processRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Print("Wrong HTTP method: " + r.Method)
		http.Error(w, `{"error": "Ошибка метода запроса. Allowed Methods: `+http.MethodPost+`"}`, http.StatusMethodNotAllowed)
		return
	}
	jwt, err := getJWtFromCookies(r)
	if err == nil {
		log.Print("You are registrated yet123:" + jwt)
		username, err := getUsernameFromCookies(r)
		if err == nil {
			log.Println(username)
		}
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
	unvalidatedForm, err := processRequestParser(r)
	if err != nil {
		log.Println(err)
		http.Error(w, `{"error": "Ошибка парсинга формы"}`, http.StatusBadGateway)
		return
	}
	user := newForm(w, unvalidatedForm)
	if err = login(w, user); err != nil {
		log.Println("Error from login after newForm " + err.Error())
		//http.Error(w, `{"error": "Ошибка регистрации: `+err.Error()+`"}`, http.StatusBadGateway)
		return
	}
	log.Println("Go to profile after login...")
	http.Redirect(w, r, "/profile", http.StatusFound)
	return
}

func processProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		log.Print("Wrong HTTP method: " + r.Method)
		http.Error(w, `{"error": "Ошибка метода запроса. Allowed Methods: `+http.MethodPut+`"}`, http.StatusMethodNotAllowed)
		return
	}
	username, err := getUsernameFromCookies(r)
	if err != nil {
		log.Print("not logged in")
		log.Print(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	unvalidatedForm, err := processRequestParser(r)
	if err != nil {
		http.Error(w, `{"error": "Ошибка парсинга формы`+err.Error()+`"}`, http.StatusBadGateway)
		return
	}

	editForm(w, username, unvalidatedForm)
	http.Redirect(w, r, "/profile", http.StatusFound)
}

func newForm(w http.ResponseWriter, unvalidatedForm *JsonForm) types.User {
	lastusername, err := database.GetLastUsername()
	var newusername string
	if err != nil {
		newusername = "FormUser_1"
	} else {
		usl := strings.Split(lastusername, "_")
		lastnum, _ := strconv.Atoi(usl[1])
		lastnumstr := strconv.Itoa(lastnum + 1)
		newusername = "FormUser_" + lastnumstr
	}

	user := types.User{}
	user.Username = newusername
	password, err := generatePassword(10)
	if err != nil {
		log.Print(err)
	}
	user.Password, err = types.HashPassword(os.Getenv("SALT") + password)

	if err != nil {
		log.Print(err)
	}

	//Здесь password нужно отправлять пользователю в ответе, причём ровно один раз
	var formerrors types.FormErrors
	var f types.Form
	err = validate(&f, unvalidatedForm, &formerrors)
	if err != nil {
		log.Print(err)

		errors_json, _ := json.Marshal(formerrors)
		//clearCookies(w)
		setErrorsCookie(w, errors_json)
	} else {
		setSuccessCookie(w)

		err := database.WriteForm(&f, &user)
		if err != nil {
			log.Print(err)
		}
		setUsernameCookie(w, newusername)
		setPasswordCookie(w, password)
		//login(w, types.User{Username: newusername, Password: password})
	}

	form_json, _ := json.Marshal(f)
	setFormDataCookie(w, form_json)
	return types.User{Username: newusername, Password: password}
}

func editForm(w http.ResponseWriter, username string, unvalidatedForm *JsonForm) {
	var formerrors types.FormErrors
	var f types.Form
	err := validate(&f, unvalidatedForm, &formerrors)
	if err != nil {
		log.Print(err)

		errors_json, _ := json.Marshal(formerrors)
		//clearCookies(w)
		setErrorsCookie(w, errors_json)
	} else {
		setSuccessCookie(w)

		err := database.UpdateForm(&f, username)
		if err != nil {
			log.Print(err)
		}
	}

	form_json, _ := json.Marshal(f)
	setFormDataCookie(w, form_json)
}

func validate(f *types.Form, form *JsonForm, formerrors *types.FormErrors) (err error) {
	var finalres bool = true
	var check bool = false
	var gen bool = false

	{
		var v string = form.Fio
		r, err := regexp.Compile(`^[A-Za-zА-Яа-яЁё\s]{1,150}$`)
		if err != nil {
			log.Print(err)
		}
		if !r.MatchString(v) {
			finalres = false
			formerrors.Fio = "Invalid fio"
			//*formerrors = append(*formerrors, 1)
		} else {
			f.Fio = v
		}
	}

	{
		var v string = form.Tel
		r, err := regexp.Compile(`^\+[0-9]{1,29}$`)
		if err != nil {
			log.Print(err)
		}
		if !r.MatchString(v) {
			finalres = false
			formerrors.Tel = "Invalid telephone"
			//*formerrors = append(*formerrors, 2)
		} else {
			f.Tel = v
		}
	}

	{
		var v string = form.Email
		r, err := regexp.Compile(`^[A-Za-z0-9._%+-]{1,30}@[A-Za-z0-9.-]{1,20}\.[A-Za-z]{1,10}$`)
		if err != nil {
			log.Print(err)
		}
		if !r.MatchString(v) {
			finalres = false
			formerrors.Email = "Invalid email"
			//*formerrors = append(*formerrors, 3)
		} else {
			f.Email = v
		}
	}

	{
		var v string = form.Date
		r, err := regexp.Compile(`^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$`)
		if err != nil {
			log.Print(err)
		}
		if !r.MatchString(v) {
			finalres = false
			formerrors.Date = "Invalid date"
			//*formerrors = append(*formerrors, 4)
		} else {
			f.Date = v
		}
	}

	{
		var v string = form.Gender
		if v != "Male" && v != "Female" {
			gen = false
		} else {
			gen = true
			f.Gender = v
		}
	}

	{
		var v string = form.Bio
		r, err := regexp.Compile(`^[A-Za-zА-Яа-яЁё;,.:0-9\-!?"'\s]{0,}$`)
		if err != nil {
			log.Print(err)
		}
		if !r.MatchString(v) {
			finalres = false
			formerrors.Bio = "Restricted symbols in bio"
		} else {
			f.Bio = v
		}
	}

	{
		var v string = form.Familiar

		if v == "on" {
			check = true
		}
	}

	{
		for _, p := range form.Favlangs {

			if p < 1 || p > 11 {
				finalres = false
				formerrors.Favlangs = "Invalid favourite langs"
				//*formerrors = append(*formerrors, 6)
				break
			} else {
				f.Favlangs = append(f.Favlangs, p)
			}
		}
	}

	if !gen {
		finalres = false
		formerrors.Gender = "Invalid gender"
		//*formerrors = append(*formerrors, 5)
	}
	if !check {
		finalres = false
		formerrors.Familiar = "Invalid familiar checkbox"
		//*formerrors = append(*formerrors, 8)
	}
	if finalres {
		return nil
	}

	return errors.New("validation failed")
}

func generatePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789__"
	password := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))
	for i := range password {
		index, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", fmt.Errorf("error generating random index: %v", err)
		}
		password[i] = charset[index.Int64()]
	}

	return string(password), nil
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	username, err := getUsernameFromCookies(r)
	if err != nil {
		lastusername, err := database.GetLastUsername()
		var newusername string
		if err != nil {
			newusername = "FormUser_1"
		} else {
			usl := strings.Split(lastusername, "_")
			lastnum, _ := strconv.Atoi(usl[1])
			lastnumstr := strconv.Itoa(lastnum + 1)
			newusername = "FormUser_" + lastnumstr
		}

		user := types.User{}
		user.Username = newusername
		password, err := generatePassword(10)
		if err != nil {
			log.Print(err)
		}
		user.Password, err = types.HashPassword(os.Getenv("SALT") + password)
		if err != nil {
			log.Print(err)
		}

		//Здесь password нужно отправлять пользователю в ответе, причём ровно один раз
		var formerrors types.FormErrors
		if err := r.ParseForm(); err != nil {
			http.Error(w, `{"error": "Ошибка парсинга формы"}`, http.StatusBadGateway)
			return
		}

		var f types.Form
		err = validate(&f, castToJsonForm(r.Form), &formerrors)
		if err != nil {
			log.Print(err)

			errors_json, _ := json.Marshal(formerrors)
			//clearCookies(w)
			setErrorsCookie(w, errors_json)
		} else {
			setSuccessCookie(w)

			err := database.WriteForm(&f, &user)
			if err != nil {
				log.Print(err)
			}
			setUsernameCookie(w, newusername)
			setPasswordCookie(w, password)
			login(w, types.User{Username: newusername, Password: password})
		}

		form_json, _ := json.Marshal(f)
		setFormDataCookie(w, form_json)
		http.Redirect(w, r, "/profile", http.StatusSeeOther)

	} else {
		var formerrors types.FormErrors
		if err := r.ParseForm(); err != nil {
			http.Error(w, `{"error": "Ошибка парсинга формы"}`, http.StatusBadGateway)
			return
		}

		var f types.Form
		err = validate(&f, castToJsonForm(r.Form), &formerrors)
		if err != nil {
			log.Print(err)

			errors_json, _ := json.Marshal(formerrors)
			//clearCookies(w)
			setErrorsCookie(w, errors_json)
		} else {
			setSuccessCookie(w)

			err := database.UpdateForm(&f, username)
			if err != nil {
				log.Print(err)
			}
		}

		form_json, _ := json.Marshal(f)
		setFormDataCookie(w, form_json)
		http.Redirect(w, r, "/form", http.StatusSeeOther)
	}
}
