package signature_test

import (
	"testing"

	"github.com/kitabisa/perkakas/v3/signature"
	"github.com/stretchr/testify/assert"
)

func TestGenerateHMACTrue(t *testing.T) {
	message := "this is my message"
	secretKey := "123-qwe"

	signature := signature.GenerateHmac(message, secretKey)
	expectedSignature := "c67f60277d9c15a90c50b4e528a24e83dfdb2adeb06e71d27a96a09b10bf754c"
	assert.Equal(t, expectedSignature, signature)
}

func TestGenerateHMACFalse(t *testing.T) {
	message := "this is my message"
	secretKey := "123-qwe"

	signature := signature.GenerateHmac(message, secretKey)
	expectedSignature := "205f52370c0a6565c60c07803211e14396f6e952a570a5fb35111766b8843a7c"
	assert.NotEqual(t, expectedSignature, signature)
}

func TestIsMatchHmacTrue(t *testing.T) {
	message := "this is my message"
	secretKey := "123-qwe"
	signatureHeader := "c67f60277d9c15a90c50b4e528a24e83dfdb2adeb06e71d27a96a09b10bf754c"

	result := signature.IsMatchHmac(message, signatureHeader, secretKey)
	expected := true
	assert.Equal(t, expected, result)
}

func TestIsMatchHmacFalse(t *testing.T) {
	message := "this is my message"
	secretKey := "123-qwe"
	signatureHeader := "205f52370c0a6565c60c07803211e14396f6e952a570a5fb35111766b8843a7c"

	result := signature.IsMatchHmac(message, signatureHeader, secretKey)
	expected := false
	assert.Equal(t, expected, result)
}
