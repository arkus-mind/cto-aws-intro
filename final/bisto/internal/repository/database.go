package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "appchalchifinal.c0s6enbo77uu.us-east-2.rds.amazonaws.com"
	port     = 5432
	user     = "chalchi"
	dbname   = "practices"
	password = "chalchi12345"
)

func CreateConnection() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	// check db
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
