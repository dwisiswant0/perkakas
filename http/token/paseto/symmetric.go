package paseto

import (
	"errors"

	"github.com/o1egl/paseto"
)

type PasetoSymmetric struct {
	Key string
}

// NewSymmetric key must be 32 bytes (32 characters long)
func NewSymmetric(key string) (*PasetoSymmetric, error) {
	// Must be 32 bytes
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes")
	}

	return &PasetoSymmetric{
		Key: key,
	}, nil
}

func (p PasetoSymmetric) Encrypt(token paseto.JSONToken, footer string) (string, error) {
	encToken, err := pasetoV2.Encrypt([]byte(p.Key), token, footer)
	if err != nil {
		return "", err
	}

	return encToken, nil
}

func (p PasetoSymmetric) Decrypt(encToken string) (paseto.JSONToken, string, error) {
	var token paseto.JSONToken
	var footer string
	err := pasetoV2.Decrypt(encToken, []byte(p.Key), &token, &footer)
	if err != nil {
		return paseto.JSONToken{}, "", err
	}

	return token, footer, nil
}
