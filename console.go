//go:build wasm && js

package gs

var Console = ConsoleType{Object: Object{Value: Global.Get("console")}}

// TODO: finish all console methods
type ConsoleType struct {
	Object
}

// Assert writes an error message to the console if the assertion is false. If
// the assertion is true, nothing happens.
//
//	c.Assert(assertion, obj1)
//	c.Assert(assertion, obj1, obj2)
//	c.Assert(assertion, obj1, obj2, /* ... ,*/ objN)

func (c ConsoleType) Assert(assertion Boolean, objs ...Valuer) {
	_, _ = c.Call("assert", append([]Valuer{assertion}, objs...)...)
}

// AssertSubst writes an error message to the console if the assertion is false.
// If the assertion is true, nothing happens.
//
//	c.AssertSubst(assertion, msg)
//	c.AssertSubst(assertion, msg, subst1)
//	c.AssertSubst(assertion, msg, subst1, /* ... ,*/ substN)
func (c ConsoleType) AssertSubst(assertion Boolean, msg String, substs ...Valuer) {
	_, _ = c.Call("assert", append([]Valuer{assertion, msg}, substs...)...)
}

// Clear clears the console if the console allows it. A graphical console, like those running on browsers, will allow it; a console displaying on the terminal, like the one running on Node, will not support it, and will have no effect (and no error).
//
//	c.Clear()
func (c ConsoleType) Clear() {
	_, _ = c.Call("clear")
}

// Count logs the number of times that this particular call to Count has been called.
//
//	c.Count() // "default: 1"
//	c.Count() // "default: 2"
func (c ConsoleType) Count() {
	_, _ = c.Call("count")
}

// CountLabel logs the number of times that this particular call to CountLabel has been called.
//
//	c.CountLabel("foo") // "foo: 1"
//	c.CountLabel("foo") // "foo: 2"
func (c ConsoleType) CountLabel(label String) {
	_, _ = c.Call("count", label)
}

// CountReset resets counter used with Count.
//
//	c.Count() // "default: 1"
//	c.Count() // "default: 2"
//	c.CountReset()
//	c.Count() // "default: 1"
func (c ConsoleType) CountReset() {
	_, _ = c.Call("reset")
}

// CountResetLabel resets counter used with CountLabel.
//
//	c.CountLabel("foo") // "foo: 1"
//	c.CountLabel("foo") // "foo: 2"
//	c.CountResetLabel("foo")
//	c.CountLabel("foo") // "foo: 1"
func (c ConsoleType) CountResetLabel(label String) {
	_, _ = c.Call("reset", label)
}

// Debug outputs a message to the web console at the "debug" log level. The message is only displayed to the user if the console is configured to display debug output. In most cases, the log level is configured within the console UI. This log level might correspond to the Debug or Verbose log level.
//
//	c.Debug(obj1)
//	c.Debug(obj1, /* ..., */ objN)
func (c ConsoleType) Debug(objs ...Valuer) {
	_, _ = c.Call("debug", objs...)
}

// DebugSubst outputs a message to the web console at the "debug" log level. The message is only displayed to the user if the console is configured to display debug output. In most cases, the log level is configured within the console UI. This log level might correspond to the Debug or Verbose log level.
//
// c.DebugSubst(msg)
// c.DebugSubst(msg, subst1, /* ..., */ substN)
func (c ConsoleType) DebugSubst(msg String, substs ...Valuer) {
	_, _ = c.Call("debug", append([]Valuer{msg}, substs...)...)
}

// Dir displays an interactive list of the properties of the specified JavaScript object. The output is presented as a hierarchical listing with disclosure triangles that let you see the contents of child objects.
// In other words, Dir is the way to see all the properties of a specified JavaScript object in console by which the developer can easily get the properties of the object.
//
//	c.Dir(object)
func (c ConsoleType) Dir(o Object) {
	_, _ = c.Call("dir", o)
}

// DirXML displays an interactive tree of the descendant elements of the specified XML/HTML element. If it is not possible to display as an element the JavaScript Object view is shown instead. The output is presented as a hierarchical listing of expandable nodes that let you see the contents of child nodes.
//
//	c.DirXML(object)
func (c ConsoleType) DirXML(o Object) {
	_, _ = c.Call("dirxml", o)
}
