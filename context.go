package v8go

// #import <stdlib.h>
// #import "v8go.h"
import "C"
import (
	"fmt"
	"io"
	"runtime"
	"unsafe"
)

type jsErr struct {
	msg      string
	location string
	stack    string
}

func (e *jsErr) Error() string {
	return e.msg
}

func (e *jsErr) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.msg)
			if e.location != "" {
				fmt.Fprintf(s, ": %s", e.location)
			}
			if e.stack != "" {
				fmt.Fprintf(s, "\n%s", e.stack)
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.msg)
	case 'q':
		fmt.Fprintf(s, "%q", e.msg)
	}
}

// Context is a global root execution environment that allows separate,
// unrelated, JavaScript applications to run in a single instance of V8.
type Context struct {
	iso *Isolate
	ptr C.ContextPtr
}

// NewContext creates a new JavaScript context for a given isoltate;
// if isolate is `nil` than a new isolate will be created.
func NewContext(iso *Isolate) (*Context, error) {
	if iso == nil {
		var err error
		iso, err = NewIsolate()
		if err != nil {
			return nil, fmt.Errorf("context: failed to create new Isolate: %v", err)
		}
	}

	// TODO: [RC] does the isolate need to track all the contexts created?
	// any script run against the context should make sure the VM has not been
	// terninated
	ctx := &Context{
		iso: iso,
		ptr: C.NewContext(iso.ptr),
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
	return &jsErr{msg: C.GoString(rtn.error), location: "blah line 8", stack: "bad things happen"}
}
