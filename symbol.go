//go:build wasm && js

package gs

var _ Valuer = Symbol{}

// TODO: fill with well-known symbols (6.1.5.1)
var (
	_ Symbol
)

type Symbol struct {
	Value
}

func (s Symbol) ValueOf() Value {
	return s.Value
}
