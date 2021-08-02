package middleware

import (
	"context"
	"testing"

	"github.com/kitabisa/perkakas/v2/token/jwt"
	"github.com/stretchr/testify/assert"
)

func TestSetClaimContext(t *testing.T) {
	claims := &jwt.UserClaim{
		UserID:         12345,
		SecondaryID:    "abcdefghijkl",
		ClientID:       "mobile",
		Scopes:         []string{"all"},
	}
	ctx := setClaimContext(context.Background(), claims)
	userID, _ := ctx.Value("UserID").(int64)
	secondaryID, _ := ctx.Value("SecondaryID").(string)
	clientID, _ := ctx.Value("ClientID").(string)
	scopes, _ := ctx.Value("Scopes").([]string)

	assert.Equal(t, int64(12345), userID)
	assert.Equal(t, "abcdefghijkl", secondaryID)
	assert.Equal(t, "mobile", clientID)
	assert.Equal(t, []string{"all"}, scopes)
}
