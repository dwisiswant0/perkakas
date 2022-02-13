package http

import (
	"context"
	"errors"
	"github.com/kitabisa/perkakas/http/token/jwt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type ConfigOpts func(s *PerkakasHttpHandler)

func NewJWT(hctx HttpHandlerContext, signKey []byte) func(next http.Handler) http.Handler {
	jwtt := jwt.NewJWT(signKey)
	writer := CustomWriter{
		C: hctx,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := bearerAuth(r, jwtt)
			if err != nil {
				log.Println(err.Error())
				writer.WriteError(w, ErrUnauthorized)
				return
			}

			parentCtx := r.Context()
			ctx := context.WithValue(parentCtx, "token", claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func bearerAuth(r *http.Request, jwtt *jwt.JWT) (*jwt.UserClaim, error) {
	authorization := r.Header.Get("Authorization")
	match, err := regexp.MatchString("^Bearer .+", authorization)
	if !match {
		return nil, errors.New("invalid jwt token")
	}
	if err != nil {
		return nil, err
	}

	tokenString := strings.Split(authorization, " ")

	token, err := jwtt.Parse(tokenString[1])
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.UserClaim)
	if !ok {
		return nil, errors.New("invalid jwt token")
	}

	return claims, nil
}
