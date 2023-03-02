//go:build wasm && js

package gs

var Global = GlobalType{
	Object: Object{Value: PredefValue(5, TypeFlagObject)},
}

// TODO: global properties/methods
type GlobalType struct {
	Object
}
