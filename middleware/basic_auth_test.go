package middleware

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuth(t *testing.T) {
	var r, _ = http.NewRequest(http.MethodGet, "/", nil)
	user := "user"
	pass := "pass"
	val := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pass)))

	r.Header.Add("Authorization", fmt.Sprintf("Basic %s", val))
	ok, err := basicAuth(r, user, pass)
	assert.Equal(t, true, ok)
	assert.Nil(t, err)
}
