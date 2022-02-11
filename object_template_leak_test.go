// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Ignore leaks within Go standard libraries http/https support code.
// The getaddrinfo detected leaks can be avoided using GODEBUG=netdns=go but
// currently there are more for loading system root certificates on macOS.
//go:build !leakcheck || !darwin
// +build !leakcheck !darwin

package v8go_test

import (
	"testing"

	v8 "rogchap.com/v8go"
)

func TestObjectTemplateLeakCheck(t *testing.T) {
	isolate := v8.NewIsolate()
	global := v8.NewObjectTemplate(isolate)

	cb := func(info *v8.FunctionCallbackInfo) *v8.Value {
		// referencing global seems to be the cause of the leak
		_ = global

		return v8.Null(isolate)
	}

	global.Set("fn", v8.NewFunctionTemplate(isolate, cb), v8.ReadOnly)

	context := v8.NewContext(isolate, global)
	context.Close()

	isolate.Dispose()
}
