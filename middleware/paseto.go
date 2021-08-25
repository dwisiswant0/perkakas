package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	phttp "github.com/kitabisa/perkakas/v2/http"
	"github.com/kitabisa/perkakas/v2/structs"
	"github.com/kitabisa/perkakas/v2/token/paseto"
	libpaseto "github.com/o1egl/paseto"
	"github.com/rs/zerolog/log"
)

func NewPaseto(hctx phttp.HttpHandlerContext, publicKey string) func(next http.Handler) http.Handler {
	pst, err := paseto.NewAsymmetric(publicKey, "")
	if err != nil {
		panic(err)
	}

	writer := phttp.CustomWriter{
		C: hctx,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, footer, err := decrypt(r, pst)
			if err != nil {
				log.Error().Msg(err.Error())
				writer.WriteError(w, structs.ErrUnauthorized)
				return
			}

			err = token.Validate()
			if err != nil {
				tokenValidationErr := fmt.Errorf("paseto token validation: %w", err)
				log.Error().Msg(tokenValidationErr.Error())
				writer.WriteError(w, structs.ErrUnauthorized)
			}

			parentCtx := r.Context()
			ctxToken := context.WithValue(parentCtx, "token", token)
			ctx := context.WithValue(ctxToken, "token_footer", footer)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func decrypt(r *http.Request, pst *paseto.PasetoAsymmetric) (libpaseto.JSONToken, string, error) {
	authorization := r.Header.Get("Authorization")
	match, err := regexp.MatchString("^Bearer .+", authorization)
	if err != nil || !match {
		return libpaseto.JSONToken{}, "", errors.New("invalid bearer token")
	}

	tokenString := strings.Split(authorization, " ")

	token, footer, err := pst.Decrypt(tokenString[1])
	if err != nil {
		return libpaseto.JSONToken{}, "", err
	}

	return token, footer, nil
}
