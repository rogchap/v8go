// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"context"
	"errors"
	"runtime"
	"unsafe"
)

// FunctionCallback is a callback that is executed in Go when a function is executed in JS.
type FunctionCallback func(info *FunctionCallbackInfo) (Valuer, error)

// FunctionCallbackInfo is the argument that is passed to a FunctionCallback.
type FunctionCallbackInfo struct {
	ctx  *ExecContext
	args []*Value
}

// Context is the current context that the callback is being executed in.
func (i *FunctionCallbackInfo) ExecContext() *ExecContext {
	return i.ctx
}

// Context will return context.
func (i *FunctionCallbackInfo) Context() context.Context {
	return i.ctx.Context()
}

// Args returns a slice of the value arguments that are passed to the JS function.
func (i *FunctionCallbackInfo) Args() []*Value {
	return i.args
}

// FunctionTemplate is used to create functions at runtime.
// There can only be one function created from a FunctionTemplate in a context.
// The lifetime of the created function is equal to the lifetime of the context.
type FunctionTemplate struct {
	*template
}

// NewFunctionTemplate creates a FunctionTemplate for a given callback.
func NewFunctionTemplate(iso *Isolate, callback FunctionCallback) (*FunctionTemplate, error) {
	if iso == nil {
		return nil, errors.New("v8go: failed to create new FunctionTemplate: Isolate cannot be <nil>")
	}
	if callback == nil {
		return nil, errors.New("v8go: failed to create new FunctionTemplate: FunctionCallback cannot be <nil>")
	}

	cbref := iso.registerCallback(callback)

	tmpl := &template{
		ptr: C.NewFunctionTemplate(iso.ptr, C.int(cbref)),
		iso: iso,
	}
	runtime.SetFinalizer(tmpl, (*template).finalizer)
	return &FunctionTemplate{tmpl}, nil
}

// GetFunction returns an instance of this function template bound to the given context.
func (tmpl *FunctionTemplate) GetFunction(ctx *ExecContext) *Function {
	val_ptr := C.FunctionTemplateGetFunction(tmpl.ptr, ctx.ptr)
	return &Function{&Value{val_ptr, ctx}}
}

// GetObject creates a new Object based on the template.
func (tmpl *FunctionTemplate) ContextValue(ctx *ExecContext) (Valuer, error) {
	return tmpl.GetFunction(ctx), nil
}

//export goFunctionCallback
func goFunctionCallback(ctxref int, cbref int, args *C.ValuePtr, argsCount int) C.ValuePtr {
	ctx := getContext(ctxref)

	info := &FunctionCallbackInfo{
		ctx:  ctx,
		args: make([]*Value, argsCount),
	}

	argv := (*[1 << 30]C.ValuePtr)(unsafe.Pointer(args))[:argsCount:argsCount]
	for i, v := range argv {
		val := &Value{ptr: v, ctx: ctx}
		info.args[i] = val
	}

	callbackFunc := ctx.iso.getCallback(cbref)
	val, err := callbackFunc(info)

	if err != nil {
		cmsg := C.CString(err.Error())
		defer C.free(unsafe.Pointer(cmsg))
		C.ThrowException(ctx.iso.ptr, cmsg)
		return nil
	}

	if val != nil && val.value() != nil {
		return val.value().ptr
	}
	return nil
}
