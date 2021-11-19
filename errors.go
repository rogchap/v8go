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
			// The StackTrace starts with the Message, so only the former needs to be printed
			io.WriteString(s, e.StackTrace)

			// If it was a compile time error, then there wouldn't be a runtime stack trace,
			// but StackTrace will still include the Message, making them equal. In this case,
			// we want to include the Location where the compilation failed.
			if e.StackTrace == e.Message && e.Location != "" {
				fmt.Fprintf(s, " (at %s)", e.Location)
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Message)
	case 'q':
		fmt.Fprintf(s, "%q", e.Message)
	}
}
