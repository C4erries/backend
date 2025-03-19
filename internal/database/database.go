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
	WriteForm(*types.Form, *types.User) error
	UpdateForm(*types.Form, string) error
	CheckUser(*types.User) error
	GetForm(string) (types.Form, error)
}

func WriteForm(f *types.Form, u *types.User) (err error) {

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

	var insertsql = []string{
		"INSERT INTO forms ",
		"(username, fio, tel, email, birth_date, gender, bio) ",
		"VALUES ($1, $2, $3, $4, $5, $6, $7)",
	}

	_, err = db.Exec(strings.Join(insertsql, ""), u.Username, f.Fio, f.Tel,
		f.Email, f.Date, f.Gender, f.Bio)
	if err != nil {
		log.Print("INSERT INTO forms ABORTED")
		return err
	}

	for _, v := range f.Favlangs {
		_, err = db.Exec("INSERT INTO favlangs VALUES ($1, $2)", u.Username, v)
		if err != nil {
			log.Println("INSERT INTO favlangs ABORTED")
			return err
		}
	}

	_, err = db.Exec("INSERT INTO userinfo VALUES ($1, $2)", u.Username, u.Password)
	if err != nil {
		log.Println("INSERT INTO userinfo ABORTED")
		return err
	}

	return nil
}

func UpdateForm(f *types.Form, username string) (err error) {
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

	var updatesql = []string{
		"UPDATE forms ",
		"SET fio = $2, tel = $3, email = $4, birth_date = $5, gender = $6, bio = $7 ",
		"WHERE username = $1",
	}

	_, err = db.Exec(strings.Join(updatesql, ""), username, f.Fio, f.Tel,
		f.Email, f.Date, f.Gender, f.Bio)
	if err != nil {
		log.Print("UPDATE forms ABORTED")
		return err
	}

	_, err = db.Exec("DELETE FROM favlangs WHERE username=$1", username)
	if err != nil {
		log.Println("DELETE FROM favlangs ABORTED")
		return err
	}

	for _, v := range f.Favlangs {
		_, err = db.Exec("INSERT INTO favlangs VALUES ($1, $2)", username, v)
		if err != nil {
			log.Println("INSERT INTO favlangs ABORTED")
			return err
		}
	}

	return nil
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

	if err := types.CheckPassword([]byte(dbpassword), u.Password); err != nil {
		return err
	}

	return nil
}

func GetForm(username string) (types.Form, error) {
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
		return types.Form{}, err
	}

	defer db.Close()

	var selectsql = []string{
		"SELECT ",
		"fio, tel, email, birth_date, gender, bio ",
		"FROM forms WHERE username=$1",
	}

	form := types.Form{}
	row := db.QueryRow(strings.Join(selectsql, ""), username)
	dateparts := strings.Split(form.Date, "T")
	form.Date = dateparts[0]
	err = row.Scan(&form.Fio, &form.Tel, &form.Email, &form.Date, &form.Gender, &form.Bio)
	if err != nil {
		return types.Form{}, err
	}

	rows, err := db.Query("SELECT lang_id FROM favlangs WHERE username=$1", username)
	if err != nil {
		return types.Form{}, err
	}

	for rows.Next() {
		var langid int
		err := rows.Scan(&langid)
		if err != nil {
			return types.Form{}, err
		}
		form.Favlangs = append(form.Favlangs, langid)
	}

	return form, nil
}

func GetLastUsername() (string, error) {
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
		return "", err
	}

	defer db.Close()

	var username string
	err = db.QueryRow("SELECT username FROM userinfo ORDER BY username DESC LIMIT 1").Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
