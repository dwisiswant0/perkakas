// Package to generate and validate paseto token. Only support paseto token v2

package paseto

import (
	"github.com/o1egl/paseto"
)

var pasetoV2 *paseto.V2

func init() {
	pasetoV2 = paseto.NewV2()
}
