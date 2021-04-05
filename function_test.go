// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	"rogchap.com/v8go"
)

func TestFunctionCall(t *testing.T) {
	t.Parallel()

	ctx, err := v8go.NewContext()
	failIf(t, err)
	_, err = ctx.RunScript("function add(a, b) { return a + b; }", "")
	failIf(t, err)
	addValue, err := ctx.Global().Get("add")
	failIf(t, err)
	iso, _ := ctx.Isolate()

	arg1, err := v8go.NewValue(iso, int32(1))
	failIf(t, err)

	fn, _ := addValue.AsFunction()
	resultValue, err := fn.Call([]v8go.Valuer{arg1, arg1})
	failIf(t, err)

	if resultValue.Int32() != 2 {
		t.Errorf("expected 1 + 1 = 2, got: %v", resultValue.DetailString())
	}
}

func TestFunctionCallError(t *testing.T) {
	t.Parallel()

	ctx, err := v8go.NewContext()
	failIf(t, err)
	_, err = ctx.RunScript("function throws() { throw 'error'; }", "script.js")
	failIf(t, err)
	addValue, err := ctx.Global().Get("throws")
	failIf(t, err)

	fn, _ := addValue.AsFunction()
	_, err = fn.Call(nil)
	if err == nil {
		t.Errorf("expected an error, got none")
	}
	got := *(err.(*v8go.JSError))
	want := v8go.JSError{Message: "error", Location: "script.js:1:21"}
	if got != want {
		t.Errorf("want %+v, got: %+v", want, got)
	}
}

func failIf(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
