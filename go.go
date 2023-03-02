//go:build wasm && js

package gs

// Go is the instance of the Go class in JavaScript
var Go = GoType{Object: Object{Value: PredefValue(6, TypeFlagObject)}}

// TODO: implement js Go type
type GoType struct {
	Object
}
