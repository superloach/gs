//go:build wasm && js

package gs

var ArrayConstructor = Function{Global.Get("Array"), 0}

// TODO: implement Array
type Array struct {
	Object
}
