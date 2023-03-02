//go:build wasm && js

package gs

import (
	"runtime"
	_ "unsafe"
)

var _ Valuer = Object{}

var ObjectConstructor = Function{Global.Get("Object"), 0}

type Object struct {
	Value

	methods map[string]Function
}

func ObjectOf(v Valuer) (Object, bool) {
	vv := v.ValueOf()

	if !vv.Type().IsObject() {
		return Object{}, false
	}

	return Object{
		Value: vv,
	}, true
}

type MethodError struct {
	Method string
}

func (m MethodError) Error() string {
	return "no such method " + m.Method
}

// Call does a JavaScript call to the method m of object o with the given
// arguments.
// It panics if v has no method m.
// The arguments get mapped to JavaScript values according to the ValueOf
// function.
func (o Object) Call(m string, args ...Valuer) (Value, error) {
	argVals, argRefs := MakeArgs(args)

	if o.methods == nil {
		o.methods = map[string]Function{}
	}

	fn, ok := o.methods[m]
	if !ok {
		fn, ok = FunctionOf(o.Get(m))
		if !ok {
			return Undefined.Value, MethodError{Method: m}

		}

		o.methods[m] = fn
	}

	res, ok := valueInvoke(fn.Ref, argRefs)
	val := MakeValue(res)

	runtime.KeepAlive(o)
	runtime.KeepAlive(argVals)

	if !ok {
		err, ok := ObjectOf(val)
		if !ok {
			panic("non-object error")
		}

		return Undefined.Value, Error{Object: err}
	}

	return val, nil
}

func (o Object) ValueOf() Value {
	return o.Value
}

func (o Object) HasOwnProperty(p String) bool {
	has, err := o.Call("hasOwnProperty", p)
	if err != nil {
		panic("hasOwnProperty should not error but: " + err.Error())
	}

	return Boolean{Value: has}.Bool()
}

// Get returns the JavaScript property p of object o.
func (o Object) Get(p string) Value {
	r := MakeValue(valueGet(o.Ref, p))
	runtime.KeepAlive(o)
	return r
}

//go:linkname valueGet syscall/js.valueGet
func valueGet(v Ref, p string) Ref
