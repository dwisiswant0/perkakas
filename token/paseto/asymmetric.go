package paseto

import (
	"crypto/ed25519"
	"encoding/hex"

	"github.com/o1egl/paseto"
)

type PasetoAsymmetric struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

func NewAsymmetric(pubKey, privKey string) (*PasetoAsymmetric, error) {
	b, err := hex.DecodeString(privKey)
	if err != nil {
		return nil, err
	}

	privateKey := ed25519.PrivateKey(b)

	b, err = hex.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}

	publicKey := ed25519.PublicKey(b)

	return &PasetoAsymmetric{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

func (p PasetoAsymmetric) Encrypt(token paseto.JSONToken, footer string) (string, error) {
	encToken, err := pasetoV2.Sign(p.PrivateKey, token, footer)
	if err != nil {
		return "", err
	}

	return encToken, nil
}

func (p PasetoAsymmetric) Decrypt(encToken string) (paseto.JSONToken, string, error) {
	var token paseto.JSONToken
	var footer string

	err := pasetoV2.Verify(encToken, p.PublicKey, &token, &footer)
	if err != nil {
		return paseto.JSONToken{}, "", err
	}

	return token, footer, nil
}
