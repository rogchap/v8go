package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

type FunctionCallback func(info *FunctionCallbackInfo) *Value

type FunctionCallbackInfo struct {
	ctx  *Context
	args []*Value
}

func (i *FunctionCallbackInfo) Context() *Context {
	return i.ctx
}

func (i *FunctionCallbackInfo) Args() []*Value {
	return i.args
}

type FunctionTemplate struct {
	*template
}

func NewFunctionTemplate(iso *Isolate, callback FunctionCallback) (*FunctionTemplate, error) {
	if iso == nil {
		return nil, errors.New("v8go: failed to create new FunctionTemplate: Isolate cannot be <nil>")
	}

	cbref := iso.registerCallback(callback)

	tmpl := &template{
		ptr: C.NewFunctionTemplate(iso.ptr, C.int(cbref)),
		iso: iso,
	}
	runtime.SetFinalizer(tmpl, (*template).finalizer)
	return &FunctionTemplate{tmpl}, nil
}

//export goFunctionCallback
func goFunctionCallback(ctxref int, cbref int, args *C.ValuePtr, args_count int) C.ValuePtr {
	ctx := getContext(ctxref)

	info := &FunctionCallbackInfo{
		ctx:  ctx,
		args: make([]*Value, args_count),
	}

	argv := (*[1 << 30]C.ValuePtr)(unsafe.Pointer(args))[:args_count:args_count]
	for i, v := range argv {
		val := &Value{ptr: v}
		runtime.SetFinalizer(val, (*Value).finalizer)
		info.args[i] = val
	}

	callbackFunc := ctx.iso.getCallback(cbref)
	if val := callbackFunc(info); val != nil {
		return val.ptr
	}
	return nil
}
