package http

import (
	"errors"
	"log"
	"net/http"
)

func NewBasicAuth(hctx HttpHandlerContext, definedUsername, definedPassword string) func(next http.Handler) http.Handler {
	writer := CustomWriter{
		C: hctx,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ok, err := basicAuth(r, definedUsername, definedPassword)
			if err != nil {
				log.Println(err.Error())
				writer.WriteError(w, ErrUnauthorized)
				return
			}

			if !ok {
				log.Println("Failed login using basic auth")
				writer.WriteError(w, ErrUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func basicAuth(r *http.Request, definedUsername, definedPassword string) (bool, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return false, errors.New("failed to parse basic auth string")
	}

	return username == definedUsername && password == definedPassword, nil
}
