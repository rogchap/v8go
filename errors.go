// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"fmt"
	"io"
	"unsafe"
)

// JSError is an error that is returned if there is are any
// JavaScript exceptions handled in the context. When used with the fmt
// verb `%+v`, will output the JavaScript stack trace, if available.
type JSError struct {
	Message    string
	Location   string
	StackTrace string
}

func newJSError(rtnErr C.RtnError) error {
	err := &JSError{
		Message:    C.GoString(rtnErr.msg),
		Location:   C.GoString(rtnErr.location),
		StackTrace: C.GoString(rtnErr.stack),
	}
	C.free(unsafe.Pointer(rtnErr.msg))
	C.free(unsafe.Pointer(rtnErr.location))
	C.free(unsafe.Pointer(rtnErr.stack))
	return err
}

func (e *JSError) Error() string {
	return e.Message
}

// Format implements the fmt.Formatter interface to provide a custom formatter
// primarily to output the javascript stack trace with %+v
func (e *JSError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') && e.StackTrace != "" {
			io.WriteString(s, e.StackTrace)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Message)
	case 'q':
		fmt.Fprintf(s, "%q", e.Message)
	}
}
