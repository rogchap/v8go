package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import "unsafe"

// JSONParse tries to parse the string and returns it as *Value if successful.
// Any JS errors will be returned as `JSError`.
func JSONParse(ctx *Context, str string) (*Value, error) {
	if ctx == nil {
		return nil, noContextErr
	}
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))

	rtn := C.JSONParse(ctx.ptr, cstr)
	return getValue(ctx, rtn), getError(rtn)
}

// JSONStringify tries to stringify the JSON-serializable object value and returns it as string.
func JSONStringify(ctx *Context, val *Value) (string, error) {
	if ctx == nil {
		return "", noContextErr
	}
	str := C.JSONStringify(ctx.ptr, val.ptr)
	defer C.free(unsafe.Pointer(str))
	return C.GoString(str), nil
}
