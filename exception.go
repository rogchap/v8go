// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

import (
	// #include <stdlib.h>
	// #include "v8go.h"
	"C"

	"fmt"
	"unsafe"
)

// NewRangeError creates a RangeError.
func NewRangeError(iso *Isolate, msg string) *Exception {
	return newExceptionError(iso, C.ERROR_RANGE, msg)
}

// NewReferenceError creates a ReferenceError.
func NewReferenceError(iso *Isolate, msg string) *Exception {
	return newExceptionError(iso, C.ERROR_REFERENCE, msg)
}

// NewSyntaxError creates a SyntaxError.
func NewSyntaxError(iso *Isolate, msg string) *Exception {
	return newExceptionError(iso, C.ERROR_SYNTAX, msg)
}

// NewTypeError creates a TypeError.
func NewTypeError(iso *Isolate, msg string) *Exception {
	return newExceptionError(iso, C.ERROR_TYPE, msg)
}

// NewWasmCompileError creates a WasmCompileError.
func NewWasmCompileError(iso *Isolate, msg string) *Exception {
	return newExceptionError(iso, C.ERROR_WASM_COMPILE, msg)
}

// NewWasmLinkError creates a WasmLinkError.
func NewWasmLinkError(iso *Isolate, msg string) *Exception {
	return newExceptionError(iso, C.ERROR_WASM_LINK, msg)
}

// NewWasmRuntimeError creates a WasmRuntimeError.
func NewWasmRuntimeError(iso *Isolate, msg string) *Exception {
	return newExceptionError(iso, C.ERROR_WASM_RUNTIME, msg)
}

// NewError creates an Error, which is the common thing to throw from
// user code.
func NewError(iso *Isolate, msg string) *Exception {
	return newExceptionError(iso, C.ERROR_GENERIC, msg)
}

func newExceptionError(iso *Isolate, typ C.ErrorTypeIndex, msg string) *Exception {
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))
	eptr := C.NewValueError(iso.ptr, typ, cmsg)
	if eptr == nil {
		panic(fmt.Errorf("invalid error type index: %d", typ))
	}
	return &Exception{&Value{ptr: eptr}}
}

// An Exception is a JavaScript exception.
type Exception struct {
	*Value
}

// value implements Valuer.
func (e *Exception) value() *Value {
	return e.Value
}

// Error implements error.
func (e *Exception) Error() string {
	return e.String()
}

// As provides support for errors.As.
func (e *Exception) As(target interface{}) bool {
	ep, ok := target.(**Exception)
	if !ok {
		return false
	}
	*ep = e
	return true
}

// Is provides support for errors.Is.
func (e *Exception) Is(err error) bool {
	eerr, ok := err.(*Exception)
	if !ok {
		return false
	}
	return eerr.String() == e.String()
}

// String implements fmt.Stringer.
func (e *Exception) String() string {
	if e.Value == nil {
		return "<nil>"
	}
	s := C.ExceptionGetMessageString(e.ptr)
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}
