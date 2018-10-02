package middlewares

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// Auth Middleware for protecting urls from unauthorized users
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		signature, err := r.Cookie("signature")
		if err != nil {
			responseJSON(http.StatusUnauthorized, w, apiError{1, "No signature cookie found"})
			return
		}

		headerPayload, err := r.Cookie("header.payload")
		if err != nil {
			responseJSON(http.StatusUnauthorized, w, apiError{1, "No headerPayload cookie found"})
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
				Expires: time.Now().AddDate(0, 0, -1),
			}
			headerPayloadCookie := http.Cookie{
				Name:    "signature",
				Expires: time.Now().AddDate(0, 0, -1),
			}
			http.SetCookie(w, &signatureCoookie)
			http.SetCookie(w, &headerPayloadCookie)
			responseJSON(http.StatusUnauthorized, w, apiError{1, err.Error()})
			return
		}

		context.Set(r, "claims", claims)

		next.ServeHTTP(w, r)
	}
}
