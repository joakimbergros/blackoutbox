package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Panicf("Unable to start server: %s", err.Error())
	}
}
