package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

// Context is a global root execution environment that allows separate,
// unrelated, JavaScript applications to run in a single instance of V8.
type Context struct {
	ptr C.ContextPtr
	iso *Isolate
}

type contextOptions struct {
	iso   *Isolate
	gTmpl *ObjectTemplate
}

// ContextOption sets options such as Isolate and Global Template to the NewContext
type ContextOption interface {
	apply(*contextOptions)
}

// NewContext creates a new JavaScript context; if no Isolate is passed as a
// ContextOption than a new Isolate will be created.
func NewContext(opt ...ContextOption) (*Context, error) {
	opts := contextOptions{}
	for _, o := range opt {
		if o != nil {
			o.apply(&opts)
		}
	}

	if opts.iso == nil {
		var err error
		opts.iso, err = NewIsolate()
		if err != nil {
			return nil, fmt.Errorf("v8go: failed to create new Isolate: %v", err)
		}
	}

	if opts.gTmpl == nil {
		opts.gTmpl = &ObjectTemplate{}
	}

	ctx := &Context{
		iso: opts.iso,
		ptr: C.NewContext(opts.iso.ptr, opts.gTmpl.ptr),
	}
	runtime.SetFinalizer(ctx, (*Context).finalizer)
	// TODO: [RC] catch any C++ exceptions and return as error
	return ctx, nil
}

// Isolate gets the current context's parent isolate.An  error is returned
// if the isolate has been terninated.
func (c *Context) Isolate() (*Isolate, error) {
	// TODO: [RC] check to see if the isolate has not been terninated
	return c.iso, nil
}

// RunScript executes the source JavaScript; origin or filename provides a
// reference for the script and used in the stack trace if there is an error.
// error will be of type `JSError` of not nil.
func (c *Context) RunScript(source string, origin string) (*Value, error) {
	cSource := C.CString(source)
	cOrigin := C.CString(origin)
	defer C.free(unsafe.Pointer(cSource))
	defer C.free(unsafe.Pointer(cOrigin))

	rtn := C.RunScript(c.ptr, cSource, cOrigin)
	return getValue(c, rtn), getError(rtn)
}

// Close will dispose the context and free the memory.
func (c *Context) Close() {
	c.finalizer()
}

func (c *Context) finalizer() {
	C.ContextDispose(c.ptr)
	c.ptr = nil
	runtime.SetFinalizer(c, nil)
}

func getValue(ctx *Context, rtn C.RtnValue) *Value {
	if rtn.value == nil {
		return nil
	}
	v := &Value{rtn.value, ctx}
	runtime.SetFinalizer(v, (*Value).finalizer)
	return v
}

func getError(rtn C.RtnValue) error {
	if rtn.error.msg == nil {
		return nil
	}
	err := &JSError{
		Message:    C.GoString(rtn.error.msg),
		Location:   C.GoString(rtn.error.location),
		StackTrace: C.GoString(rtn.error.stack),
	}
	C.free(unsafe.Pointer(rtn.error.msg))
	C.free(unsafe.Pointer(rtn.error.location))
	C.free(unsafe.Pointer(rtn.error.stack))
	return err
}
