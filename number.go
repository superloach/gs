//go:build wasm && js

package gs

var _ Valuer = Number{}

// TODO: implement number operations (6.1.6.1)
type Number struct {
	Object
}

func (n Number) ValueOf() Value {
	return n.Value
}
