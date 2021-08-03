package middleware

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strings"

	phttp "github.com/kitabisa/perkakas/v2/http"
	"github.com/kitabisa/perkakas/v2/structs"
	"github.com/kitabisa/perkakas/v2/token/jwt"
	"github.com/rs/zerolog/log"
)

func NewJWT(hctx phttp.HttpHandlerContext, signKey []byte) func(next http.Handler) http.Handler {
	jwtt := jwt.NewJWT(signKey)
	writer := phttp.CustomWriter{
		C: hctx,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := bearerAuth(r, jwtt)
			if err != nil {
				log.Error().Msg(err.Error())
				writer.WriteError(w, structs.ErrUnauthorized)
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
	if err != nil || !match {
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
