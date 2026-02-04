package middleware

import (
	"log"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO implement authentication
		log.Println("Authenticating...")

		next.ServeHTTP(w, r)
	})
}
