package middlewares

import (
	"net/http"
)

// Post Middleware for Post connection
func Post(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			responseJSON(http.StatusMethodNotAllowed, w, apiError{1, "Method Not Allowed"})
			return
		}
		next.ServeHTTP(w, r)
	}
}

// Get Middleware for Get connection
func Get(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			responseJSON(http.StatusMethodNotAllowed, w, apiError{1, "Method Not Allowed"})
			return
		}
		next.ServeHTTP(w, r)
	}
}

// Put Middleware for Put connection
func Put(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			responseJSON(http.StatusMethodNotAllowed, w, apiError{1, "Method Not Allowed"})
			return
		}
		next.ServeHTTP(w, r)
	}
}

// Delete Middleware for Delete connection
func Delete(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			responseJSON(http.StatusMethodNotAllowed, w, apiError{1, "Method Not Allowed"})
			return
		}
		next.ServeHTTP(w, r)
	}
}

// GetPut Middleware for Get or Put connection
func GetPut(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPut {
			responseJSON(http.StatusMethodNotAllowed, w, apiError{1, "Method Not Allowed"})
			return
		}
		next.ServeHTTP(w, r)
	}
}
