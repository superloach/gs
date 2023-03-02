//go:build wasm && js

package gs

var _ Valuer = BigInt{}

// TODO: implement BigInt operations (6.1.6.2)
type BigInt struct {
	Value
}

// BigIntOf converts a JavaScript value into a BigInt, if possible.
//
// TODO: coerce non-bigint values into bigints?
func BigIntOf(v Valuer) (BigInt, bool) {
	vv := v.ValueOf()

	if vv.Type() != TypeBigInt {
		return BigInt{}, false
	}

	return BigInt{
		Value: vv,
	}, true
}

func (b BigInt) ValueOf() Value {
	return b.Value
}
