// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

// Error creates a generic error message.
func Error(iso *Isolate, msg string) *Value {
	fmt.Printf("iso = %+v\n", iso)
	cstr := C.CString(msg)
	defer C.free(unsafe.Pointer(cstr))
	ptr := C.ExceptionError(iso.ptr, cstr)
	v := &Value{ptr: ptr}
	runtime.SetFinalizer(v, (*Value).finalizer)
	return v
}

// RangeError creates a range error message.
func RangeError(iso *Isolate, msg string) *Value {
	cstr := C.CString(msg)
	defer C.free(unsafe.Pointer(cstr))
	ptr := C.ExceptionRangeError(iso.ptr, cstr)
	v := &Value{ptr: ptr}
	runtime.SetFinalizer(v, (*Value).finalizer)
	return v
}

// ReferenceError creates a reference error message.
func ReferenceError(iso *Isolate, msg string) *Value {
	cstr := C.CString(msg)
	defer C.free(unsafe.Pointer(cstr))
	ptr := C.ExceptionReferenceError(iso.ptr, cstr)
	v := &Value{ptr: ptr}
	runtime.SetFinalizer(v, (*Value).finalizer)
	return v
}

// SyntaxError creates a syntax error message.
func SyntaxError(iso *Isolate, msg string) *Value {
	cstr := C.CString(msg)
	defer C.free(unsafe.Pointer(cstr))
	ptr := C.ExceptionSyntaxError(iso.ptr, cstr)
	v := &Value{ptr: ptr}
	runtime.SetFinalizer(v, (*Value).finalizer)
	return v
}

// TypeError creates a type error message.
func TypeError(iso *Isolate, msg string) *Value {
	cstr := C.CString(msg)
	defer C.free(unsafe.Pointer(cstr))
	ptr := C.ExceptionTypeError(iso.ptr, cstr)
	v := &Value{ptr: ptr}
	runtime.SetFinalizer(v, (*Value).finalizer)
	return v
}
