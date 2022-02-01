package random

import (
	"math/rand"

	"github.com/google/uuid"
)

func UUID() (id string) {
	return uuid.New().String()
}

type IRandomStringGenerator interface {
	String(length int) string
}

type RandomStringGenerator struct{}

func NewRandomStringGenerator() *RandomStringGenerator {
	return &RandomStringGenerator{}
}

// String random string with n length
func (s *RandomStringGenerator) String(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
