//go:build wasm && js

package gs

var Null = NullType{Value: PredefValue(2, TypeFlagNone)}

type NullType struct {
	Value
}

func (n NullType) ValueOf() Value {
	return n.Value
}
