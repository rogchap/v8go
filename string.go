// Copyright 2022 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"unsafe"
)

// String is a JavaScript string object (ECMA-262, 4.3.17)
type String struct {
	*Value
}

func jsStringResult(ctx *Context, rtn C.RtnValue) (String, error) {
	if rtn.value == nil {
		return String{nil}, newJSError(rtn.error)
	}
	return String{&Value{rtn.value, ctx}}, nil
}

// NewString returns a JS string from the Go string or an error if the string is longer than the max V8 string length.
func NewString(iso *Isolate, val string) (String, error) {
	cstr := C.CString(val)
	defer C.free(unsafe.Pointer(cstr))
	rtn := C.NewValueString(iso.ptr, cstr, C.int(len(val)))
	return jsStringResult(nil, rtn)
}

// MustNewString wraps NewString with a panic on error check.
//
// Use for strings known to be within than the max V8 string length.
// V8's max string length is (1 << 28) - 16 on 32-bit systems
// and (1 << 29) - 24 on other systems, at the time of writing.
func MustNewString(iso *Isolate, val string) String {
	jsStr, err := NewString(iso, val)
	if err != nil {
		panic(err)
	}
	return jsStr
}
