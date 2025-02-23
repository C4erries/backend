package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"os"

	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type form struct {
	fio      string
	tel      string
	email    string
	date     string
	gender   string
	favlangs []int
	bio      string
}

func process(w http.ResponseWriter, r *http.Request) {
	var formerrors []int

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET POST PUT OPTIONS CONNECT HEAD")
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, `{"error": "Ошибка парсинга формы"}`, http.StatusBadGateway)
		return
	}
	var f form
	err := Validate(&f, r.Form, &formerrors)
	if err != nil {
		log.Print(err)
		res, _ := json.Marshal(formerrors)
		w.WriteHeader(400)
		w.Write(res)

	} else {
		err := WriteForm(&f)
		if err != nil {
			log.Print(err)
		}
		w.WriteHeader(200)

	}

}

func Validate(f *form, form url.Values, formerrors *[]int) (err error) {
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
				f.fio = v
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
				f.tel = v
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
				f.email = v
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
				f.date = v
			}
		}

		if key == "Gender" {
			var v string = value[0]
			if v != "Male" && v != "Female" {
				gen = false
			} else {
				gen = true
				f.gender = v
			}
		}

		if key == "Bio" {
			var v string = value[0]
			f.bio = v
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
						f.favlangs = append(f.favlangs, np)
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

func WriteForm(f *form) (err error) {

	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")
	/*
		postgresHost := os.Getenv("POSTGRES_HOST")
		connectStr := "host=" + postgresHost + " user=" + postgresUser +
		" password=" + postgresPassword +
		" dbname=" + postgresDB + " sslmode=disable"
	*/
	//postgresUser := "postgres"
	//postgresPassword := "123"
	//postgresDB := "back3"
	connectStr := "user=" + postgresUser +
		" password=" + postgresPassword +
		" dbname=" + postgresDB + " sslmode=disable"
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		return err
	}
	defer db.Close()
	var insertsql = []string{
		"INSERT INTO forms",
		"(fio, tel, email, birth_date, gender, bio)",
		"VALUES ($1, $2, $3, $4, $5, $6) returning form_id",
	}
	var form_id int
	err = db.QueryRow(strings.Join(insertsql, ""), f.fio, f.tel,
		f.email, f.date, f.gender, f.bio).Scan(&form_id)
	if err != nil {
		log.Print("YEP")
		return err
	}

	for _, v := range f.favlangs {
		_, err = db.Exec("INSERT INTO favlangs VALUES ($1, $2)", form_id, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var port string = os.Getenv("APP_PORT")
	server := http.Server{
		Addr: "0.0.0.0:" + port,
	}
	http.HandleFunc("/process", process)
	log.Println("starting server..")
	server.ListenAndServe()
}
