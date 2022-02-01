package token

import (
	"crypto/sha256"
	"fmt"
)

// Simple token is token using sum 256 and a salt.
func SimpleToken(data, salt string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(data+salt)))
}
