package jwt

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt"
)

type JWTRSA struct {
	pubKey    []byte
	privKey   []byte
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
}

func NewJWTRSA(pubKey []byte, privKey []byte) (jwtrsa *JWTRSA, err error) {
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(privKey)
	if err != nil {
		return nil, err
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		return nil, err
	}

	jwtrsa = &JWTRSA{
		pubKey:    pubKey,
		privKey:   privKey,
		signKey:   signKey,
		verifyKey: verifyKey,
	}

	return
}

func (j *JWTRSA) Create(claims jwt.Claims) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err = token.SignedString(j.signKey)
	return
}

func (j *JWTRSA) Parse(token string) (tkn *jwt.Token, err error) {
	tkn, err = jwt.ParseWithClaims(token, &UserClaim{}, j.keyFunction)
	return
}

func (j *JWTRSA) keyFunction(token *jwt.Token) (interface{}, error) {
	return j.verifyKey, nil
}
