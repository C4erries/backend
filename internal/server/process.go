package server

import (
	"backend/internal/database"
	"backend/internal/types"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

func ProcessHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var formerrors []int

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, `{"error": "Ошибка парсинга формы"}`, http.StatusBadGateway)
		return
	}

	var f types.Form
	err := validate(&f, r.Form, &formerrors)
	if err != nil {
		log.Print(err)

		errors_json, _ := json.Marshal(formerrors)
		setErrorsCookie(w, errors_json)
	} else {
		setSucsessCookie(w)

		err := database.WriteForm(ctx, &f)
		if err != nil {
			log.Print(err)
		}
	}

	form_json, _ := json.Marshal(f)
	setFormDataCookie(w, form_json)
}

func validate(f *types.Form, form url.Values, formerrors *[]int) (err error) {
	var check bool = false
	var gen bool = false
	for key, value := range form {

		if key == "Fio" {
			var v string = value[0]
			r, err := regexp.Compile(`^[A-Za-zА-Яа-яЁё\s]{1,150}$`)
			if err != nil {
				log.Print(err)
			}
			if !r.MatchString(v) {
				*formerrors = append(*formerrors, 1)
			} else {
				f.Fio = v
			}
		}

		if key == "Tel" {
			var v string = value[0]
			r, err := regexp.Compile(`^\+[0-9]{1,29}$`)
			if err != nil {
				log.Print(err)
			}
			if !r.MatchString(v) {
				*formerrors = append(*formerrors, 2)
			} else {
				f.Tel = v
			}
		}

		if key == "Email" {
			var v string = value[0]
			r, err := regexp.Compile(`^[A-Za-z0-9._%+-]{1,30}@[A-Za-z0-9.-]{1,20}\.[A-Za-z]{1,10}$`)
			if err != nil {
				log.Print(err)
			}
			if !r.MatchString(v) {
				*formerrors = append(*formerrors, 3)
			} else {
				f.Email = v
			}
		}

		if key == "Birth_date" {
			var v string = value[0]
			r, err := regexp.Compile(`^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$`)
			if err != nil {
				log.Print(err)
			}
			if !r.MatchString(v) {
				*formerrors = append(*formerrors, 4)
			} else {
				f.Date = v
			}
		}

		if key == "Gender" {
			var v string = value[0]
			if v != "Male" && v != "Female" {
				gen = false
			} else {
				gen = true
				f.Gender = v
			}
		}

		if key == "Bio" {
			var v string = value[0]
			f.Bio = v
		}

		if key == "Familiar" {
			var v string = value[0]

			if v == "on" {
				check = true
			}
		}

		if key == "Favlangs" {
			for _, p := range value {
				np, err := strconv.Atoi(p)
				if err != nil {
					log.Print(err)
					*formerrors = append(*formerrors, 6)
					break
				} else {
					if np < 1 || np > 11 {
						*formerrors = append(*formerrors, 6)
						break
					} else {
						f.Favlangs = append(f.Favlangs, np)
					}
				}
			}
		}
	}
	if !gen {
		*formerrors = append(*formerrors, 5)
	}
	if !check {
		*formerrors = append(*formerrors, 8)
	}
	if len(*formerrors) == 0 {
		return nil
	}

	return errors.New("validation failed")
}
