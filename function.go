// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wasm && js

package gs

import (
	"fmt"
	"runtime"
	"sync"
	_ "unsafe"
)

var (
	//go:linkname funcsMu syscall/js.funcsMu
	funcsMu sync.Mutex

	//go:linkname funcs syscall/js.funcs
	funcs map[uint32]func(Value, []Value) any

	//go:linkname nextFuncID syscall/js.nextFuncID
	nextFuncID uint32
)

// Function is a wrapped Go function to be called by JavaScript.
type Function struct {
	Value // the JavaScript function that invokes the Go function
	id    uint32
}

// TODO: document functionof
func FunctionOf(v Valuer) (Function, bool) {
	vv := v.ValueOf()

	if vv.Type() != TypeFunction {
		return Function{}, false
	}

	return Function{
		Value: vv,
	}, true
}

// WrapFunction returns a function to be used by JavaScript.
//
// The Go function fn is called with the value of JavaScript's "this" keyword and the
// arguments of the invocation. The return value of the invocation is
// the result of the Go function mapped back to JavaScript according to ValueOf.
//
// Invoking the wrapped Go function from JavaScript will
// pause the event loop and spawn a new goroutine.
// Other wrapped functions which are triggered during a call from Go to JavaScript
// get executed on the same goroutine.
//
// As a consequence, if one wrapped function blocks, JavaScript's event loop
// is blocked until that function returns. Hence, calling any async JavaScript
// API, which requires the event loop, like fetch (http.Client), will cause an
// immediate deadlock. Therefore a blocking function should explicitly start a
// new goroutine.
//
// Func.Release must be called to free up resources when the function will not be invoked any more.
func WrapFunction(fn func(this Value, args []Value) any) (Function, error) {
	funcsMu.Lock()
	id := nextFuncID
	nextFuncID++
	funcs[id] = fn
	funcsMu.Unlock()

	wrap, err := Go.Call("_makeFuncWrapper", ValueOf(id))
	if err != nil {
		return Function{}, fmt.Errorf("make func wrapper: %w", err)
	}

	return Function{
		id:    id,
		Value: wrap,
	}, nil
}

// Release frees up resources allocated for the function.
// The function must not be invoked after calling Release.
// It is allowed to call Release while the function is still running.
func (f Function) Release() {
	funcsMu.Lock()
	delete(funcs, f.id)
	funcsMu.Unlock()
}

// Invoke does a JavaScript call of the function f with the given arguments.
// The arguments get mapped to JavaScript values according to the ValueOf function.
func (f Function) Invoke(args ...Valuer) (Value, error) {
	argVals, argRefs := MakeArgs(args)
	res, ok := valueInvoke(f.Ref, argRefs)
	val := MakeValue(res)

	runtime.KeepAlive(f)
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

//go:linkname valueInvoke syscall/js.valueInvoke
func valueInvoke(v Ref, args []Ref) (Ref, bool)

// New uses JavaScript's "new" operator with value v as constructor and the given arguments.
// It panics if v is not a JavaScript function.
// The arguments get mapped to JavaScript values according to the ValueOf function.
func (f Function) New(args ...Valuer) (Value, error) {
	argVals, argRefs := MakeArgs(args)
	res, ok := valueNew(f.Ref, argRefs)
	val := MakeValue(res)
	runtime.KeepAlive(f)
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

//go:linkname valueNew syscall/js.valueNew
func valueNew(v Ref, args []Ref) (Ref, bool)

// setEventHandler is defined in the runtime package.
//
//go:linkname setEventHandler syscall/js.setEventHandler
func setEventHandler(fn func())

func init() {
	setEventHandler(handleEvent)
}

var releasedErrMsg = ToString("call to released function")

func handleEvent() {
	cb, ok := ObjectOf(Go.Get("_pendingEvent"))
	if !ok {
		return
	}
	Go.Set("_pendingEvent", Null)

	id := uint32(cb.Get("id").Int())
	if id == 0 { // zero indicates deadlock
		select {}
	}
	funcsMu.Lock()
	f, ok := funcs[id]
	funcsMu.Unlock()
	if !ok {
		_, _ = Console.Call("error", releasedErrMsg)
		return
	}

	this := cb.Get("this")
	argsObj := cb.Get("args")
	args := make([]Value, argsObj.Length())
	for i := range args {
		args[i] = argsObj.Index(i)
	}
	result := f(this, args)
	cb.Set("result", result)
}
