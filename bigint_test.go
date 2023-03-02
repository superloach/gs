package gs_test

import (
	"testing"

	"github.com/superloach/gs"
)

func TestBigInt(t *testing.T) {
	zeroBigInt := gs.ToString("0n")

	val, err := gs.Global.Call("eval", zeroBigInt)
	if err != nil {
		t.Fatalf("eval bigint: %v", err)
	}

	if ty := val.Type(); ty != gs.TypeBigInt {
		t.Fatalf("expected bigint, got %v", ty)
	}
}
