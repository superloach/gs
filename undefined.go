//go:build wasm && js

package gs

var Undefined = UndefinedType{Value: Value{Ref: 0}}

type UndefinedType struct {
	Value
}

func (u UndefinedType) ValueOf() Value {
	return u.Value
}
