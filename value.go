package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"runtime"
	"unsafe"
)

// Value represents all Javascript values and objects
type Value struct {
	ptr C.ValuePtr
}

// String will return the string representation of the value. Primitive values
// are returned as-is, objects will return `[object Object]` and functions will
// print their definition.
func (v *Value) String() string {
	s := C.ValueToString(v.ptr)
	value := C.GoString(s)
	C.free(unsafe.Pointer(s))
	return value
}

func (v *Value) finalizer() {
	C.ValueDispose(v.ptr)
	v.ptr = nil
	runtime.SetFinalizer(v, nil)
}
