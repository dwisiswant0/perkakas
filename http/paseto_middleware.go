package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/kitabisa/perkakas/v3/http/token/paseto"
	"log"
	"net/http"
	"regexp"
	"strings"

	libpaseto "github.com/o1egl/paseto"
)

func NewPaseto(hctx HttpHandlerContext, publicKey string) func(next http.Handler) http.Handler {
	pst, err := paseto.NewAsymmetric(publicKey, "")
	if err != nil {
		panic(err)
	}

	writer := CustomWriter{
		C: hctx,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, footer, err := decrypt(r, pst)
			if err != nil {
				log.Println(err.Error())
				writer.WriteError(w, ErrUnauthorized)
				return
			}

			err = token.Validate()
			if err != nil {
				tokenValidationErr := fmt.Errorf("paseto token validation: %w", err)
				log.Println(tokenValidationErr.Error())
				writer.WriteError(w, ErrUnauthorized)
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
