// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include "v8go.h"
import "C"
import (
	"unsafe"
)

// Function is a JavaScript function.
type Function struct {
	*Value
}

// Call this JavaScript function with the given arguments.
func (fn *Function) Call(args ...Valuer) (*Value, error) {
	var argptr *C.ValuePtr
	if len(args) > 0 {
		var cArgs = make([]C.ValuePtr, len(args))
		for i, arg := range args {
			cArgs[i] = arg.value().ptr
		}
		argptr = (*C.ValuePtr)(unsafe.Pointer(&cArgs[0]))
	}
	fn.ctx.register()
	rtn := C.FunctionCall(fn.ptr, C.int(len(args)), argptr)
	fn.ctx.deregister()
	return getValue(fn.ctx, rtn), getError(rtn)
}

// Invoke a constructor function to create an object instance.
func (fn *Function) NewInstance(args ...Valuer) (*Object, error) {
	var argptr *C.ValuePtr
	if len(args) > 0 {
		var cArgs = make([]C.ValuePtr, len(args))
		for i, arg := range args {
			cArgs[i] = arg.value().ptr
		}
		argptr = (*C.ValuePtr)(unsafe.Pointer(&cArgs[0]))
	}
	fn.ctx.register()
	rtn := C.FunctionNewInstance(fn.ptr, C.int(len(args)), argptr)
	fn.ctx.deregister()
	return getObject(fn.ctx, rtn), getError(rtn)
}

// Return the script origin for a function.
func (fn *Function) SourceMapUrl() *Value {
	ptr := C.FunctionSourceMapUrl(fn.ptr)
	return &Value{ptr, fn.ctx}
}
