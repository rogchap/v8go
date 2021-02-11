// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"

// Error creates a generic error message.
func Error(msg string) *Value {
	panic("not implemented")
}

// RangeError creates a range error message.
func RangeError(msg string) *Value {
	panic("not implemented")
}

// ReferenceError creates a reference error message.
func ReferenceError(msg string) *Value {
	panic("not implemented")
}

// SyntaxError creates a syntax error message.
func SyntaxError(msg string) *Value {
	panic("not implemented")
}

// TypeError creates a type error message.
func TypeError(msg string) *Value {
	panic("not implemented")
}
