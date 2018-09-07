package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// Auth Middleware for protecting urls from unauthorized users
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		header := r.Header.Get("Authorization")
		if len(header) == 0 {
			http.Error(w, "Not authorized", 401)
			return
		}
		splittedHeader := strings.Fields(r.Header.Get("Authorization"))
		if len(splittedHeader) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}
		tokenString := splittedHeader[1]

		if len(tokenString) == 0 {
			http.Error(w, "Not authorized", 401)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				http.Error(w, "Unexpected signing method: ", 401)
				// return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte("Secret"), nil
		})

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			http.Error(w, err.Error(), 401)
			return
		} else {
			fmt.Println(claims["foo"], claims["nbf"])
		}

		next.ServeHTTP(w, r)
	}
}
