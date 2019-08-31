package v8go

// #import "v8go.h"
import "C"
import "runtime"

type Value struct {
	ptr C.ValuePtr
}

func (v *Value) String() string {
	return C.GoString(C.ValueToString(v.ptr))
}

func (v *Value) finalizer() {
	C.ValueDispose(v.ptr)
	v.ptr = nil
	runtime.SetFinalizer(v, nil)
}
