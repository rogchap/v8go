// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// RegisterCallback is exported for testing only.
func (i *Isolate) RegisterCallback(cb FunctionCallbackWithError) int {
	return i.registerCallback(cb)
}

// GetCallback is exported for testing only.
func (i *Isolate) GetCallback(ref int) FunctionCallbackWithError {
	return i.getCallback(ref)
}

// GetContext is exported for testing only.
var GetContext = getContext

// Ref is exported for testing only.
func (c *Context) Ref() int {
	return c.ref
}
