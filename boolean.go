//go:build wasm && js

package gs

var (
	True  = Boolean{Value: PredefValue(3, TypeFlagNone)}
	False = Boolean{Value: PredefValue(4, TypeFlagNone)}
)

type Boolean struct {
	Value
}

func ToBoolean(b bool) Boolean {
	if b {
		return True
	}

	return False
}

func (b Boolean) ValueOf() Value {
	return b.Value
}

// Bool returns the boolean b as a bool.
// It panics if b is not true or false.
func (b Boolean) Bool() bool {
	switch b.Ref {
	case True.Ref:
		return true
	case False.Ref:
		return false
	default:
		panic("bool not true or false")
	}
}
