//go:build wasm && js

package gs_test

import (
	"testing"

	"github.com/superloach/gs"
)

func TestGlobalThis(t *testing.T) {
	_ = gs.Global
}
