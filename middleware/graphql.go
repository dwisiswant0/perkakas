package middleware

import (
	"context"
	"net/http"
)

// MapHeaderToContext is used for assigning header to the request context to be processed in graphql resolver.
func MapHeaderToContext(next http.Handler) (wrapped http.Handler) {
	wrapped = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key := range r.Header {
			r = r.WithContext(context.WithValue(r.Context(), key, r.Header.Get(key)))
		}
		next.ServeHTTP(w, r)
	})
	return
}
