package gs

var Reflect = ReflectType{
	Object: Object{
		Value: Global.Get("Reflect"),
	},
}

// TODO: finish Reflect methods and docs
type ReflectType struct {
	Object
}

func (r ReflectType) Constructor(v Valuer) Object {
	con, err := r.Call("constructor", v.ValueOf())
	if err != nil {
		panic("constructor: " + err.Error())
	}

	o, ok := ObjectOf(con)
	if !ok {
		panic("not object")
	}

	return o
}

func (r ReflectType) Get(o Object, p String) Value {
	v, err := r.Call("get", o, p)
	if err != nil {
		panic("get: " + err.Error())
	}

	return v
}
