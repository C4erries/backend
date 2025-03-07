package database

import (
	"backend/internal/types"
	"context"
	"database/sql"
	"log"
	"os"
	"strings"
	"time"
)

type Database interface {
	DatabaseInit()
	WriteForm(*types.Form) error
}

var singleFlight (*SingleFlight)

func WriteForm(parentCtx context.Context, f *types.Form) (err error) {
	current_ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
	defer cancel()
	_, err = singleFlight.Do(current_ctx, f.Email, func(ctx context.Context) (any, error) {
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
			return nil, err
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
			return nil, err
		}

		for _, v := range f.Favlangs {
			_, err = db.Exec("INSERT INTO favlangs VALUES ($1, $2)", form_id, v)
			if err != nil {
				log.Println("INSERT INTO favlangs aborted")
				return nil, err
			}
		}
		return nil, nil
	})
	return err

}

func DatabaseInit() {
	singleFlight = NewSingleFlight()
}
