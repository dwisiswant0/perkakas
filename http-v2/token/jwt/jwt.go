package jwt

import (
	"github.com/golang-jwt/jwt"
)

type JWT struct {
	SignKey []byte
}

// UserClaim defines user token claim
type UserClaim struct {
	UserID      int64    `json:"user_id"`
	SecondaryID string   `json:"secondary_id"`
	ClientID    string   `json:"client_id"`
	Scopes      []string `json:"scopes"`
	jwt.StandardClaims
}

func NewJWT(signKey []byte) *JWT {
	return &JWT{
		SignKey: signKey,
	}
}

func (j *JWT) Create(claims jwt.Claims) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(j.SignKey)
	return
}

func (j *JWT) Parse(token string) (tkn *jwt.Token, err error) {
	tkn, err = jwt.ParseWithClaims(token, &UserClaim{}, j.keyFunction)
	return
}

func (j *JWT) keyFunction(token *jwt.Token) (interface{}, error) {
	return []byte(j.SignKey), nil
}

func (j *JWT) Valid(token string) (valid bool, err error) {
	tkn, err := j.Parse(token)
	if err != nil {
		return
	}

	valid = tkn.Valid
	return
}
