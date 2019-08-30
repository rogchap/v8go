package v8go

// #import "v8go.h"
import "C"
import "runtime"

type Context struct {
	ptr C.ContextPtr
}

func NewContext(iso *Isolate) *Context {
	ctx := &Context{C.NewContext(iso.ptr)}
	runtime.SetFinalizer(ctx, (*Context).release)
	return ctx
}

func (c *Context) RunScript() {
	C.RunScript(c.ptr, C.CString("source"), C.CString("origin"))
}

func (c *Context) release() {
	//TODO dispose of object in C++A
	c.ptr = nil
	runtime.SetFinalizer(c, nil)
}
