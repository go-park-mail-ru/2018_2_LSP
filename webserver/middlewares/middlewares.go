package middlewares

import (
	"net/http"
)

// Middleware Http middleware
type Middleware func(next http.HandlerFunc) http.HandlerFunc

// Chain Util for chaining different middlewares into new one
func Chain(mw ...Middleware) Middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last(w, r)
		}
	}
}
