package jwt

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	signKey = []byte("abcde")
)

func TestCreate(t *testing.T) {
	claims := UserClaim{}
	claims.UserID = 12345
	claims.SecondaryID = "8fae85be-e441-4344-8634-d41f23684146"
	claims.Scopes = []string{"read"}
	claims.ClientID = "apdifuoqpweyr9823u"
	claims.Id = "63410cd1-110b-4a2c-8c3f-ae1535eda9a1"
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(48 * time.Hour).Unix()

	jwt := NewJWT(signKey)
	token, err := jwt.Create(claims)
	assert.Nil(t, err, "Error must be nil")
	t.Log(token)
}

func TestCreateRSA(t *testing.T) {
	claims := UserClaim{}
	claims.UserID = 12345
	claims.SecondaryID = "8fae85be-e441-4344-8634-d41f23684146"
	claims.Scopes = []string{"read"}
	claims.ClientID = "apdifuoqpweyr9823u"
	claims.Id = "63410cd1-110b-4a2c-8c3f-ae1535eda9a1"
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()

	pubKey, _ := ioutil.ReadFile("dummy.rsa.pub")
	privKey, _ := ioutil.ReadFile("dummy.rsa")

	jwt, _ := NewJWTRSA(pubKey, privKey)

	token, err := jwt.Create(claims)
	assert.Nil(t, err, "Failed to create token")
	t.Log(token)
}

func TestParseRSA(t *testing.T) {
	claims := UserClaim{}
	claims.UserID = 12345
	claims.SecondaryID = "8fae85be-e441-4344-8634-d41f23684146"
	claims.Scopes = []string{"read"}
	claims.ClientID = "apdifuoqpweyr9823u"
	claims.Id = "63410cd1-110b-4a2c-8c3f-ae1535eda9a1"
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()

	pubKey, _ := ioutil.ReadFile("dummy.rsa.pub")
	privKey, _ := ioutil.ReadFile("dummy.rsa")

	jwt, _ := NewJWTRSA(pubKey, privKey)

	token, err := jwt.Create(claims)
	assert.Nil(t, err, "Failed to create token")

	parsedToken, err := jwt.Parse(token)
	assert.Nil(t, err, err)
	t.Log(parsedToken.Header)
	t.Log(parsedToken.Claims)
	t.Log(parsedToken.Signature)
}

func TestParse(t *testing.T) {
	claims := UserClaim{}
	claims.UserID = 12345
	claims.SecondaryID = "8fae85be-e441-4344-8634-d41f23684146"
	claims.Scopes = []string{"read"}
	claims.ClientID = "apdifuoqpweyr9823u"
	claims.Id = "63410cd1-110b-4a2c-8c3f-ae1535eda9a1"
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Minute).Unix()

	jwt := NewJWT(signKey)
	token, err := jwt.Create(claims)
	assert.Nil(t, err, "Error must be nil 1")

	parsedToken, err := jwt.Parse(token)
	assert.Nil(t, err, err)
	t.Log(parsedToken.Header)
	t.Log(parsedToken.Claims)
	t.Log(parsedToken.Signature)
}
