package http

import (
	"context"
	"github.com/kitabisa/perkakas/http/token/jwt"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type AuthOption struct {
	SignKey  []byte // jwt sign key
	Username string // Basic Auth username
	Password string // Basuc Auth password
}

// Middleware authentication supports jwt or basic auth
func NewAuthentication(hctx HttpHandlerContext, authOption AuthOption) func(next http.Handler) http.Handler {
	jwtt := jwt.NewJWT(authOption.SignKey)
	definedUsername := authOption.Username
	definedPassword := authOption.Password
	writer := CustomWriter{
		C: hctx,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			auth := strings.ToLower(r.Header.Get("Authorization"))
			if strings.HasPrefix(auth, "basic") {
				ok, err := basicAuth(r, definedUsername, definedPassword)
				if err != nil {
					log.Println(err)
					writer.WriteError(w, ErrUnauthorized)
					return
				}

				if !ok {
					log.Println("Failed login using basic auth")
					writer.WriteError(w, ErrUnauthorized)
				}
			} else if strings.HasPrefix(auth, "bearer") {
				claims, err := bearerAuth(r, jwtt)
				if err != nil {
					log.Println(err.Error())
					writer.WriteError(w, ErrUnauthorized)
					return
				}

				ctx = setClaimContext(ctx, claims)
				ctx = context.WithValue(ctx, "token", claims) // compatibility with existing logic in all our services
			} else {
				log.Println("invalid authentication type")
				writer.WriteError(w, ErrUnauthorized)
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
