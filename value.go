// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wasm && js

// Package js gives access to the WebAssembly host environment when using the js/wasm architecture.
// Its API is based on JavaScript semantics.
//
// This package is EXPERIMENTAL. Its current scope is only to allow tests to run, but not yet to provide a
// comprehensive API for users. It is exempt from the Go compatibility promise.
package gs

import (
	"runtime"
	"strconv"
	"unsafe"
)

// Ref is used to identify a JavaScript value, since the value itself can not be passed to WebAssembly.
//
// The JavaScript value "undefined" is represented by the value 0.
// A JavaScript number (64-bit float, except 0 and NaN) is represented by its IEEE 754 binary representation.
// All other values are represented as an IEEE 754 binary representation of NaN with bits 0-31 used as
// an ID and bits 32-34 used to differentiate between string, symbol, function and object.
type Ref uint64

// NaNHead are the upper 32 bits of a ref which are set if the value is not encoded as an IEEE 754 number (see above).
const NaNHead = 0x7FF80000

type Valuer interface {
	ValueOf() Value
}

// Value represents a JavaScript value. The zero value is the JavaScript value "undefined".
// Values can be checked for equality with the Equal method.
type Value struct {
	_     [0]func() // uncomparable; to make == not compile
	Ref   Ref       // identifies a JavaScript value, see ref type
	GCPtr *Ref      // used to trigger the finalizer when the Value is not referenced any more
}

func (v Value) ValueOf() Value {
	return v
}

const (
	// the type flags need to be in sync with wasm_exec.js
	TypeFlagNone = iota
	TypeFlagObject
	TypeFlagString
	TypeFlagSymbol
	TypeFlagFunction
)

func MakeValue(r Ref) Value {
	var gcPtr *Ref
	typeFlag := (r >> 32) & 7
	if (r>>32)&NaNHead == NaNHead && typeFlag != TypeFlagNone {
		gcPtr = new(Ref)
		*gcPtr = r
		runtime.SetFinalizer(gcPtr, func(p *Ref) {
			finalizeRef(*p)
		})
	}

	return Value{Ref: r, GCPtr: gcPtr}
}

// Finalize is a wrapper for syscall/js.finalizeRef(r)
func (r Ref) Finalize() {
	finalizeRef(r)
}

//go:linkname finalizeRef syscall/js.finalizeRef
func finalizeRef(r Ref)

func PredefValue(id uint32, typeFlag byte) Value {
	return Value{Ref: (NaNHead|Ref(typeFlag))<<32 | Ref(id)}
}

func FloatValue(f float64) Value {
	if f == 0 {
		return ValueZero
	}
	if f != f {
		return ValueNaN
	}
	return Value{Ref: *(*Ref)(unsafe.Pointer(&f))}
}

// Error wraps a JavaScript error.
type Error struct {
	// Object is the underlying JavaScript error object.
	Object
}

// Error implements the error interface.
func (e Error) Error() string {
	return "JavaScript error: " + e.Get("message").String()
}

var (
	ValueNaN  = PredefValue(0, TypeFlagNone)
	ValueZero = PredefValue(1, TypeFlagNone)
)

// Equal reports whether v and w are equal according to JavaScript's === operator.
func (v Value) Equal(w Value) bool {
	return v.Ref == w.Ref && v.Ref != ValueNaN.Ref
}

// IsUndefined reports whether v is the JavaScript value "undefined".
func (v Value) IsUndefined() bool {
	return v.Ref == Undefined.Ref
}

// IsNull reports whether v is the JavaScript value "null".
func (v Value) IsNull() bool {
	return v.Ref == Null.Ref
}

// IsNaN reports whether v is the JavaScript value "NaN".
func (v Value) IsNaN() bool {
	return v.Ref == ValueNaN.Ref
}

// ValueOf returns x as a JavaScript value:
//
//	| Go                     | JavaScript             |
//	| ---------------------- | ---------------------- |
//	| js.Value               | [its value]            |
//	| js.Func                | function               |
//	| nil                    | null                   |
//	| bool                   | boolean                |
//	| integers and floats    | number                 |
//	| string                 | string                 |
//	| []interface{}          | new array              |
//	| map[string]interface{} | new object             |
//
// Panics if x is not one of the expected types.
func ValueOf(x any) Value {
	switch x := x.(type) {
	case Value:
		return x
	case Function:
		return x.Value
	case nil:
		return Null.Value
	case bool:
		if x {
			return True.Value
		} else {
			return False.Value
		}
	case int:
		return FloatValue(float64(x))
	case int8:
		return FloatValue(float64(x))
	case int16:
		return FloatValue(float64(x))
	case int32:
		return FloatValue(float64(x))
	case int64:
		return FloatValue(float64(x))
	case uint:
		return FloatValue(float64(x))
	case uint8:
		return FloatValue(float64(x))
	case uint16:
		return FloatValue(float64(x))
	case uint32:
		return FloatValue(float64(x))
	case uint64:
		return FloatValue(float64(x))
	case uintptr:
		return FloatValue(float64(x))
	case unsafe.Pointer:
		return FloatValue(float64(uintptr(x)))
	case float32:
		return FloatValue(float64(x))
	case float64:
		return FloatValue(x)
	case string:
		return MakeValue(stringVal(x))
	case []any:
		l := FloatValue(float64(len(x)))
		a, err := ArrayConstructor.New(l)
		if err != nil {
			panic("array construction error: " + err.Error())
		}

		for i, s := range x {
			a.SetIndex(i, s)
		}

		return a
	case map[string]any:
		o, err := ObjectConstructor.New()
		if err != nil {
			panic("object construction error: " + err.Error())
		}

		for k, v := range x {
			o.Set(k, v)
		}

		return o
	case Valuer:
		return x.ValueOf()
	default:
		panic("ValueOf: invalid value")
	}
}

// Type represents the JavaScript type of a Value.
type Type int

const (
	TypeUndefined Type = iota
	TypeNull
	TypeBoolean
	TypeNumber
	TypeString
	TypeSymbol
	TypeObject
	TypeFunction
	TypeBigInt
)

func (t Type) String() string {
	switch t {
	case TypeUndefined:
		return "undefined"
	case TypeNull:
		return "null"
	case TypeBoolean:
		return "boolean"
	case TypeNumber:
		return "number"
	case TypeString:
		return "string"
	case TypeSymbol:
		return "symbol"
	case TypeObject:
		return "object"
	case TypeFunction:
		return "function"
	case TypeBigInt:
		return "bigint"
	default:
		panic("bad type")
	}
}

func (t Type) IsObject() bool {
	return t == TypeObject || t == TypeFunction
}

var constructorProperty = ToString("constructor")

// Type returns the JavaScript type of the value v. It is similar to JavaScript's typeof operator,
// except that it returns TypeNull instead of TypeObject for null.
func (v Value) Type() Type {
	switch v.Ref {
	case Undefined.Ref:
		return TypeUndefined
	case Null.Ref:
		return TypeNull
	case True.Ref, False.Ref:
		return TypeBoolean
	}
	if v.IsNumber() {
		return TypeNumber
	}
	typeFlag := (v.Ref >> 32) & 7
	switch typeFlag {
	case TypeFlagObject:
		return TypeObject
	case TypeFlagString:
		return TypeString
	case TypeFlagSymbol:
		return TypeSymbol
	case TypeFlagFunction:
		return TypeFunction
	default:
		constructor := Reflect.Get(
			Reflect.Constructor(v),
			constructorProperty,
		)

		if constructor.Equal(Global.Get("BigInt")) {
			return TypeBigInt
		}

		panic("invalid type: " + strconv.Itoa(int(typeFlag)))
	}
}

// Set sets the JavaScript property p of value v to ValueOf(x).
// It panics if v is not a JavaScript object.
func (v Value) Set(p string, x any) {
	if vType := v.Type(); !vType.IsObject() {
		panic(&ValueError{"Value.Set", vType})
	}
	xv := ValueOf(x)
	valueSet(v.Ref, p, xv.Ref)
	runtime.KeepAlive(v)
	runtime.KeepAlive(xv)
}

//go:linkname valueSet syscall/js.valueSet
func valueSet(v Ref, p string, x Ref)

// Delete deletes the JavaScript property p of value v.
// It panics if v is not a JavaScript object.
func (v Value) Delete(p string) {
	if vType := v.Type(); !vType.IsObject() {
		panic(&ValueError{"Value.Delete", vType})
	}
	valueDelete(v.Ref, p)
	runtime.KeepAlive(v)
}

//go:linkname valueDelete syscall/js.valueDelete
func valueDelete(v Ref, p string)

// Index returns JavaScript index i of value v.
// It panics if v is not a JavaScript object.
func (v Value) Index(i int) Value {
	if vType := v.Type(); !vType.IsObject() {
		panic(&ValueError{"Value.Index", vType})
	}
	r := MakeValue(valueIndex(v.Ref, i))
	runtime.KeepAlive(v)
	return r
}

//go:linkname valueIndex syscall/js.valueIndex
func valueIndex(v Ref, i int) Ref

// SetIndex sets the JavaScript index i of value v to ValueOf(x).
// It panics if v is not a JavaScript object.
func (v Value) SetIndex(i int, x any) {
	if vType := v.Type(); !vType.IsObject() {
		panic(&ValueError{"Value.SetIndex", vType})
	}
	xv := ValueOf(x)
	valueSetIndex(v.Ref, i, xv.Ref)
	runtime.KeepAlive(v)
	runtime.KeepAlive(xv)
}

//go:linkname valueSetIndex syscall/js.valueSetIndex
func valueSetIndex(v Ref, i int, x Ref)

func MakeArgs(args []Valuer) ([]Value, []Ref) {
	argVals := make([]Value, len(args))
	argRefs := make([]Ref, len(args))
	for i, arg := range args {
		v := arg.ValueOf()
		argVals[i] = v
		argRefs[i] = v.Ref
	}
	return argVals, argRefs
}

// Length returns the JavaScript property "length" of v.
// It panics if v is not a JavaScript object.
func (v Value) Length() int {
	if vType := v.Type(); !vType.IsObject() {
		panic(&ValueError{"Value.SetIndex", vType})
	}
	r := valueLength(v.Ref)
	runtime.KeepAlive(v)
	return r
}

//go:linkname valueLength syscall/js.valueLength
func valueLength(v Ref) int

func (v Value) IsNumber() bool {
	return v.Ref == ValueZero.Ref ||
		v.Ref == ValueNaN.Ref ||
		(v.Ref != Undefined.Ref && (v.Ref>>32)&NaNHead != NaNHead)
}

func (v Value) float(method string) float64 {
	if !v.IsNumber() {
		panic(&ValueError{method, v.Type()})
	}
	if v.Ref == ValueZero.Ref {
		return 0
	}
	return *(*float64)(unsafe.Pointer(&v.Ref))
}

// Float returns the value v as a float64.
// It panics if v is not a JavaScript number.
func (v Value) Float() float64 {
	return v.float("Value.Float")
}

// Int returns the value v truncated to an int.
// It panics if v is not a JavaScript number.
func (v Value) Int() int {
	return int(v.float("Value.Int"))
}

// Truthy returns the JavaScript "truthiness" of the value v. In JavaScript,
// false, 0, "", null, undefined, and NaN are "falsy", and everything else is
// "truthy". See https://developer.mozilla.org/en-US/docs/Glossary/Truthy.
func (v Value) Truthy() bool {
	switch v.Type() {
	case TypeUndefined, TypeNull:
		return false
	case TypeBoolean:
		return Boolean{Value: v}.Bool()
	case TypeNumber:
		return v.Ref != ValueNaN.Ref && v.Ref != ValueZero.Ref
	case TypeString:
		return v.String() != ""
	case TypeSymbol, TypeFunction, TypeObject:
		return true
	default:
		panic("bad type")
	}
}

// String returns the value v as a string.
// String is a special case because of Go's String method convention. Unlike the other getters,
// it does not panic if v's Type is not TypeString. Instead, it returns a string of the form "<T>"
// or "<T: V>" where T is v's type and V is a string representation of v's value.
func (v Value) String() string {
	switch v.Type() {
	case TypeString:
		return jsString(v)
	case TypeUndefined:
		return "<undefined>"
	case TypeNull:
		return "<null>"
	case TypeBoolean:
		return "<boolean: " + jsString(v) + ">"
	case TypeNumber:
		return "<number: " + jsString(v) + ">"
	case TypeSymbol:
		return "<symbol>"
	case TypeObject:
		return "<object>"
	case TypeFunction:
		return "<function>"
	default:
		panic("bad type")
	}
}

func jsString(v Value) string {
	str, length := valuePrepareString(v.Ref)
	runtime.KeepAlive(v)
	b := make([]byte, length)
	valueLoadString(str, b)
	finalizeRef(str)
	return string(b)
}

//go:linkname valuePrepareString syscall/js.valuePrepareString
func valuePrepareString(v Ref) (Ref, int)

//go:linkname valueLoadString syscall/js.valueLoadString
func valueLoadString(v Ref, b []byte)

// InstanceOf reports whether v is an instance of type t according to JavaScript's instanceof operator.
func (v Value) InstanceOf(t Value) bool {
	r := valueInstanceOf(v.Ref, t.Ref)
	runtime.KeepAlive(v)
	runtime.KeepAlive(t)
	return r
}

//go:linkname valueInstanceOf syscall/js.valueInstanceOf
func valueInstanceOf(v Ref, t Ref) bool

// A ValueError occurs when a Value method is invoked on
// a Value that does not support it. Such cases are documented
// in the description of each method.
type ValueError struct {
	Method string
	Type   Type
}

func (e *ValueError) Error() string {
	return "syscall/js: call of " + e.Method + " on " + e.Type.String()
}
