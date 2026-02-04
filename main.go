package main

import (
	"blackoutbox/internal/handlers/documents"
	"blackoutbox/internal/stores"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "file:./app.db?_foreign_keys=on")
	if err != nil {
		log.Panic("unable to connect to db")
	}

	documentStore := stores.DocumentStore{Db: db}
	documentHandler := documents.DocumentHandler{Store: &documentStore}

	mux := http.NewServeMux()

	mux.Handle("GET /documents", documentHandler.Get())
	mux.Handle("GET /documents/{id}", documentHandler.GetById())
	mux.Handle("POST /documents", documentHandler.Post())
	mux.Handle("PATCH /documents", documentHandler.Update())

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Panicf("Unable to start server: %s", err.Error())
	}
}
