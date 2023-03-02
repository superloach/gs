//go:build wasm && js

package gs

import _ "unsafe"

var StringConstructor = Function{Value: Global.Get("String")}

type String struct {
	Object
}

func ToString(s string) String {
	ref := stringVal(s)
	return String{
		Object: Object{
			Value: MakeValue(ref),
		},
	}
}

//go:linkname stringVal syscall/js.stringVal
func stringVal(x string) Ref

func (s String) ValueOf() Value {
	return s.Value
}

func (s String) IndexOf(search String) int {
	idx, _ := s.Call("indexOf", search)
	return idx.Int()
}
