package v8go

// #import <stdlib.h>
// #import "v8go.h"
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

type Context struct {
	ptr C.ContextPtr
}

func NewContext(iso *Isolate) *Context {
	ctx := &Context{C.NewContext(iso.ptr)}
	runtime.SetFinalizer(ctx, (*Context).finalizer)
	return ctx
}

// RunScript executes the source JavaScript; origin or filename  provides a
// reference for the script and used in the exception stack trace.
func (c *Context) RunScript(source string, origin string) (*Value, error) {
	cSource := C.CString(source)
	cOrigin := C.CString(origin)
	defer C.free(unsafe.Pointer(cSource))
	defer C.free(unsafe.Pointer(cOrigin))

	rtn := C.RunScript(c.ptr, cSource, cOrigin)
	return getValue(rtn), getError(rtn)
}

func (c *Context) finalizer() {
	C.ContextDispose(c.ptr)
	c.ptr = nil
	runtime.SetFinalizer(c, nil)
}

func getValue(rtn C.RtnValue) *Value {
	if rtn.value == nil {
		return nil
	}
	v := &Value{rtn.value}
	runtime.SetFinalizer(v, (*Value).finalizer)
	return v
}

func getError(rtn C.RtnValue) error {
	if rtn.error == nil {
		return nil
	}
	defer C.free(unsafe.Pointer(rtn.error))
	return errors.New(C.GoString(rtn.error))
}
