package middleware

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	phttp "github.com/kitabisa/perkakas/v2/http"
	"github.com/kitabisa/perkakas/v2/structs"
	"github.com/kitabisa/perkakas/v2/token/jwt"
	"github.com/rs/zerolog/log"
)

type AuthOption struct {
	SignKey  []byte // jwt sign key
	Username string // Basic Auth username
	Password string // Basuc Auth password
}

// Func authenticate supports jwt or basic auth
func Authenticate(hctx phttp.HttpHandlerContext, authOption AuthOption) func(next http.Handler) http.Handler {
	jwtt := jwt.NewJWT(authOption.SignKey)
	definedUsername := authOption.Username
	definedPassword := authOption.Password
	writer := phttp.CustomWriter{
		C: hctx,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			auth := strings.ToLower(r.Header.Get("Authorization"))
			if strings.HasPrefix(auth, "basic") {
				ok, err := basicAuth(r, definedUsername, definedPassword)
				if err != nil {
					log.Error().Msg(err.Error())
					writer.WriteError(w, structs.ErrUnauthorized)
					return
				}

				if !ok {
					log.Error().Msg("Failed login using basic auth")
					writer.WriteError(w, structs.ErrUnauthorized)
				}
			} else if strings.HasPrefix(auth, "bearer") {
				claims, err := bearerAuth(r, jwtt)
				if err != nil {
					log.Error().Msg(err.Error())
					writer.WriteError(w, structs.ErrUnauthorized)
					return
				}

				ctx = setClaimContext(ctx, claims)
				ctx = context.WithValue(ctx, "token", claims) // compatibility with existing logic in all our services
			} else {
				log.Error().Msg("invalid authentication type")
				writer.WriteError(w, structs.ErrUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func setClaimContext(ctx context.Context, claims *jwt.UserClaim) context.Context {
	e := reflect.ValueOf(claims).Elem()
	for i := 0; i < e.NumField(); i++ {
		name := e.Type().Field(i).Name
		value := e.Field(i).Interface()

		ctx = context.WithValue(ctx, name, value)
	}

	return ctx
}
