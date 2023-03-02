//go:build wasm && js

package gs

// TODO: how to access properties?
type Property struct {
	Value    Value
	Writable Boolean
	Get      Function
}
