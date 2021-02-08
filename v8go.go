// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

/*
Package v8go provides an API to execute JavaScript.
*/
package v8go

// #include "v8go.h"
import "C"

// Version returns the version of the V8 Engine with the -v8go suffix
func Version() string {
	return C.GoString(C.Version())
}
