//go:build wasm && js

package gs

import (
	"runtime"
	_ "unsafe"
)

var Uint8ArrayConstructor = Function{Value: Global.Get("Uint8Array")}

type Uint8Array struct {
	Object
}

func Uint8ArrayOf(v Valuer) (Uint8Array, bool) {
	o, ok := ObjectOf(v)
	if !ok {
		return Uint8Array{}, false
	}

	if !o.InstanceOf(Uint8ArrayConstructor.Value) {
		return Uint8Array{}, false
	}

	return Uint8Array{
		Object: o,
	}, true
}

// CopyBytesToGo copies bytes from src to dst.
// It returns the number of bytes copied, which will be the minimum of the lengths of src and dst.
func (src Uint8Array) CopyBytesToGo(dst []byte) int {
	n, _ := copyBytesToGo(dst, src.Ref)
	runtime.KeepAlive(src)
	return n
}

//go:linkname copyBytesToGo syscall/js.copyBytesToGo
func copyBytesToGo(dst []byte, src Ref) (int, bool)

// CopyBytesToJS copies bytes from src to dst.
// It returns the number of bytes copied, which will be the minimum of the lengths of src and dst.
func (dst Uint8Array) CopyBytesToJS(src []byte) int {
	n, _ := copyBytesToJS(dst.Ref, src)
	runtime.KeepAlive(dst)
	return n
}

//go:linkname copyBytesToJS syscall/js.copyBytesToJS
func copyBytesToJS(dst Ref, src []byte) (int, bool)
