package middlewares

import (
	"log"
	"net/http"
)

// Tracing Middleware for requests tracing
func Tracing(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Tracing request for %s", r.RequestURI)
		next.ServeHTTP(w, r)
	}
}
