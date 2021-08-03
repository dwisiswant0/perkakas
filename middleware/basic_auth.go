package middleware

import (
	"errors"
	"net/http"

	phttp "github.com/kitabisa/perkakas/v2/http"
	"github.com/kitabisa/perkakas/v2/log"
	"github.com/kitabisa/perkakas/v2/structs"
)

func BasicAuth(hctx phttp.HttpHandlerContext, definedUsername, definedPassword string) func(next http.Handler) http.Handler {
	writer := phttp.CustomWriter{
		C: hctx,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctxName := "Middleware.BasicAuth"
			log := log.GetSublogger(ctx, ctxName)

			ok, err := basicAuth(r, definedUsername, definedPassword)
			if err != nil {
				log.Error().Msg(err.Error())
				writer.WriteError(w, structs.ErrUnauthorized)
			}

			if !ok {
				log.Error().Msg("Failed login using basic auth")
				writer.WriteError(w, structs.ErrUnauthorized)
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
