package middlewares

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func responseJSON(statusCode int, w http.ResponseWriter, p interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(p)
}

type apiError struct {
	Code    int
	Message string
}

// Auth Middleware for protecting urls from unauthorized users
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		signature, err := r.Cookie("signature")
		if err != nil {
			responseJSON(http.StatusUnauthorized, w, "")
			return
		}

		headerPayload, err := r.Cookie("header.payload")
		if err != nil {
			responseJSON(http.StatusUnauthorized, w, "")
			return
		}

		tokenString := headerPayload.Value + "." + signature.Value
		claims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("HeAdfasdf3ref&^%$Dfrtgauyhia"), nil
		})

		if err != nil {
			signatureCoookie := http.Cookie{
				Name:    "signature",
				Expires: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
			}
			headerPayloadCookie := http.Cookie{
				Name:    "signature",
				Expires: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
			}
			http.SetCookie(w, &signatureCoookie)
			http.SetCookie(w, &headerPayloadCookie)
			responseJSON(http.StatusUnauthorized, w, apiError{1, err.Error()})
			return
		}

		next.ServeHTTP(w, r)
	}
}
