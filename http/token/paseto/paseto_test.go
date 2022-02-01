package paseto

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"
	"time"

	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/assert"
)

func TestSymmetricEncrypt(t *testing.T) {
	symmetricPaseto, err := NewSymmetric("PrU5AbXJawKJIUOJFmd4f6ZwmifLvvoF")
	assert.NoError(t, err)

	now := time.Now()

	token := paseto.JSONToken{
		Audience:   "Kitabisa services",
		Issuer:     "Kulonuwun",
		Jti:        "706cbfce-c031-4a44-815e-030f963f7d4e",
		Subject:    "cac2ee7e-70d0-4220-badd-7b5695f53ad8",
		Expiration: now.Add(24 * time.Hour),
		IssuedAt:   now,
		NotBefore:  now,
	}

	token.Set("email", "budi@kitabisa.com")
	token.Set("name", "Budi Ariyanto")

	_, err = symmetricPaseto.Encrypt(token, "Kitabisa.com")
	assert.NoError(t, err)
}

func TestSymmetricDecrypt(t *testing.T) {
	symmetricPaseto, _ := NewSymmetric("PrU5AbXJawKJIUOJFmd4f6ZwmifLvvoF")

	encToken := "v2.local.vx3Xd8jVUVJlsYwkQzEQ_XqjbuE3ejO49QnkBpL3imOGRoaKS-sXGZcJCfZekI1beLjzlYdRqnVKYDl_rC1Wp1rFbRym-a_HDM7V1Lk7xoksoxEwTOBEA4C7bX463odkc3Vge0rOudWsAOx8vQ9rV8LZtBGYqR4NVhPxdE4BdUytRXuUpHf09udPaL5gTmNGIpFFkg4MqAkqkWApyxla1Z-_mRWIe5h7wH_NE3l432s4aGnHk0iI3W3V7yVsMDRrmCDyndJI27jf1FIJ-6tRcn9KFA6s-48ZHVU4xzVjad0U6wnU61Ycz3ld3fNbz7_UjjP1bpBIzqmUAaCUPrvhaJuI0_CNqQWxdqxS1LBCUehfsebh-KfrPQUEpNbB7nOU-gAqYon8P9jyd4l9fIS6NJa4edXrv12cvkPx1BmGWzHbLtXoSX0ETA.S2l0YWJpc2EuY29t"
	decToken, footer, err := symmetricPaseto.Decrypt(encToken)
	if assert.NoError(t, err) {
		assert.Equal(t, "Kitabisa.com", footer)
		assert.Equal(t, "Kitabisa services", decToken.Audience)
		assert.Equal(t, "Kulonuwun", decToken.Issuer)
		assert.Equal(t, "706cbfce-c031-4a44-815e-030f963f7d4e", decToken.Jti)
		assert.Equal(t, "cac2ee7e-70d0-4220-badd-7b5695f53ad8", decToken.Subject)
		assert.Equal(t, "budi@kitabisa.com", decToken.Get("email"))
		assert.Equal(t, "Budi Ariyanto", decToken.Get("name"))
	}
}

func TestAsymmetricEncrypt(t *testing.T) {
	pubkey, privKey, _ := ed25519.GenerateKey(nil)
	pub := hex.EncodeToString(pubkey)
	priv := hex.EncodeToString(privKey)

	ast, err := NewAsymmetric(pub, priv)
	assert.NoError(t, err)

	now := time.Now()

	token := paseto.JSONToken{
		Audience:   "Kitabisa services",
		Issuer:     "Kulonuwun",
		Jti:        "706cbfce-c031-4a44-815e-030f963f7d4e",
		Subject:    "cac2ee7e-70d0-4220-badd-7b5695f53ad8",
		Expiration: now.Add(24 * time.Hour),
		IssuedAt:   now,
		NotBefore:  now,
	}

	token.Set("email", "budi@kitabisa.com")
	token.Set("name", "Budi Ariyanto")

	_, err = ast.Encrypt(token, "Kitabisa.com")
	assert.NoError(t, err)
}

func TestAsymmetricDecrypt(t *testing.T) {
	pub := "07ddbb390ad3096b51138358c8e082686e8023853616768d0cd60cdb5ab68b58"
	priv := "211eec1bbe7d6e4779542fb14591f8e2c0ef57607bff958b32267f1b9667bf0007ddbb390ad3096b51138358c8e082686e8023853616768d0cd60cdb5ab68b58"

	ast, _ := NewAsymmetric(pub, priv)

	encToken := "v2.public.eyJhdWQiOiJLaXRhYmlzYSBzZXJ2aWNlcyIsImVtYWlsIjoiYnVkaUBraXRhYmlzYS5jb20iLCJleHAiOiIyMDIxLTA4LTIxVDE2OjMyOjE3KzA3OjAwIiwiaWF0IjoiMjAyMS0wOC0yMFQxNjozMjoxNyswNzowMCIsImlzcyI6Ikt1bG9udXd1biIsImp0aSI6IjcwNmNiZmNlLWMwMzEtNGE0NC04MTVlLTAzMGY5NjNmN2Q0ZSIsIm5hbWUiOiJCdWRpIEFyaXlhbnRvIiwibmJmIjoiMjAyMS0wOC0yMFQxNjozMjoxNyswNzowMCIsInN1YiI6ImNhYzJlZTdlLTcwZDAtNDIyMC1iYWRkLTdiNTY5NWY1M2FkOCJ9gGMdg3eH8wFcIZ5tn3_87e7Z7jwfW-Q6vZqmkXijRTRGCFb171Vl4Rc8PRvXirvV6dj9Ns8xTBY0ksmRskz9DA.S2l0YWJpc2EuY29t"
	decToken, footer, err := ast.Decrypt(encToken)
	if assert.NoError(t, err) {
		assert.Equal(t, "Kitabisa.com", footer)
		assert.Equal(t, "Kitabisa services", decToken.Audience)
		assert.Equal(t, "Kulonuwun", decToken.Issuer)
		assert.Equal(t, "706cbfce-c031-4a44-815e-030f963f7d4e", decToken.Jti)
		assert.Equal(t, "cac2ee7e-70d0-4220-badd-7b5695f53ad8", decToken.Subject)
		assert.Equal(t, "budi@kitabisa.com", decToken.Get("email"))
		assert.Equal(t, "Budi Ariyanto", decToken.Get("name"))
	}
}
