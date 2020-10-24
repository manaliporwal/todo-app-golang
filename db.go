package main

import (
	"database/sql"
	"fmt"
)

var db *sql.DB
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "welcome"
	dbname   = "to_do_app"
	)

// hook up to postgres db
func initDB() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	    "password=%s dbname=%s sslmode=disable",
	        host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
}