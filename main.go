package main

import (
	"database/sql"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("sqlite3", "sql://app.db")
	if err != nil {
		log.Panic("unable to connect to db")
	}

	mux := http.NewServeMux()

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Panicf("Unable to start server: %s", err.Error())
	}
}
