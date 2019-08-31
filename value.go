package v8go

// #import "v8go.h"
import "C"

type Value struct {
	ptr C.ValuePtr
}

func (v *Value) String() string {
	return C.GoString(C.ValueToString(v.ptr))
}
