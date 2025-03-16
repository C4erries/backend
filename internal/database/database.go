package database

import (
	"backend/internal/types"
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"
)

type Database interface {
	WriteForm(*types.Form) error
	UpdateForm(*types.Form, string) error
	GetForm(string) (*types.Form, error)
	CheckUser(*types.User) error
}

func GetForm(string) (*types.Form, error) {
	return &types.Form{}, nil
}

func CheckUser(u *types.User) (err error) {
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")

	postgresHost := os.Getenv("POSTGRES_HOST")

	/*
		postgresHost := "db"
		postgresUser := "postgres"
		postgresPassword := "****"
		postgresDB := "back3"
	*/
	connectStr := "host=" + postgresHost + " user=" + postgresUser +
		" password=" + postgresPassword +
		" dbname=" + postgresDB + " sslmode=disable"
	//log.Println(connectStr)
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		return err
	}

	defer db.Close()

	var dbpassword string
	err = db.QueryRow("SELECT enc_password FROM userinfo WHERE username=$1", u.Username).Scan(&dbpassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("username not found")
		}
		return err
	}

	if dbpassword != u.Password {
		return errors.New("password mismatch")
	}

	return nil
}

func WriteForm(f *types.Form) (err error) {
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")

	postgresHost := os.Getenv("POSTGRES_HOST")

	connectStr := "host=" + postgresHost + " user=" + postgresUser +
		" password=" + postgresPassword +
		" dbname=" + postgresDB + " sslmode=disable"

	//log.Println(connectStr)

	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		return err
	}
	defer db.Close()
	var insertsql = []string{
		"INSERT INTO forms",
		"(fio, Tel, email, birth_date, gender, bio)",
		"VALUES ($1, $2, $3, $4, $5, $6) returning form_id",
	}
	var form_id int
	err = db.QueryRow(strings.Join(insertsql, ""), f.Fio, f.Tel,
		f.Email, f.Date, f.Gender, f.Bio).Scan(&form_id)
	if err != nil {
		log.Print("YEP")
		return err
	}

	for _, v := range f.Favlangs {
		_, err = db.Exec("INSERT INTO favlangs VALUES ($1, $2)", form_id, v)
		if err != nil {
			log.Println("INSERT INTO favlangs aborted")
			return err
		}
	}
	return nil
}
