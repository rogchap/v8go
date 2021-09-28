// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	v8 "github.com/airplanedev/v8go"
)

func TestFunctionCall(t *testing.T) {
	t.Parallel()

	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	_, err := ctx.RunScript("function add(a, b) { return a + b; }", "")
	failIf(t, err)
	addValue, err := ctx.Global().Get("add")
	failIf(t, err)
	iso := ctx.Isolate()

	arg1, err := v8.NewValue(iso, int32(1))
	failIf(t, err)

	fn, _ := addValue.AsFunction()
	resultValue, err := fn.Call(v8.Undefined(iso), arg1, arg1)
	failIf(t, err)

	if resultValue.Int32() != 2 {
		t.Errorf("expected 1 + 1 = 2, got: %v", resultValue.DetailString())
	}
}

func TestFunctionCallToGoFunc(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	defer iso.Dispose()
	global := v8.NewObjectTemplate(iso)

	called := false
	printfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		called = true
		return nil
	})

	err := global.Set("print", printfn, v8.ReadOnly)
	failIf(t, err)

	ctx := v8.NewContext(iso, global)
	defer ctx.Close()

	val, err := ctx.RunScript(`(a, b) => { print("foo"); }`, "")
	failIf(t, err)
	fn, err := val.AsFunction()
	failIf(t, err)
	resultValue, err := fn.Call(v8.Undefined(iso))
	failIf(t, err)

	if !called {
		t.Errorf("expected my function to be called, wasn't")
	}
	if !resultValue.IsUndefined() {
		t.Errorf("expected undefined, got: %v", resultValue.DetailString())
	}
}

func TestFunctionCallWithObjectReceiver(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	global := v8.NewObjectTemplate(iso)

	ctx := v8.NewContext(iso, global)
	val, err := ctx.RunScript(`class Obj { constructor(input) { this.input = input } print() { return this.input.toString() } }; new Obj("some val")`, "")
	failIf(t, err)
	obj, err := val.AsObject()
	failIf(t, err)
	fnVal, err := obj.Get("print")
	failIf(t, err)
	fn, err := fnVal.AsFunction()
	failIf(t, err)
	resultValue, err := fn.Call(obj)
	failIf(t, err)

	if !resultValue.IsString() || resultValue.String() != "some val" {
		t.Errorf("expected 'some val', got: %v", resultValue.DetailString())
	}
}

func TestFunctionCallError(t *testing.T) {
	t.Parallel()

	ctx := v8.NewContext()
	iso := ctx.Isolate()
	defer iso.Dispose()
	defer ctx.Close()

	_, err := ctx.RunScript("function throws() { throw 'error'; }", "script.js")
	failIf(t, err)
	addValue, err := ctx.Global().Get("throws")
	failIf(t, err)

	fn, _ := addValue.AsFunction()
	_, err = fn.Call(v8.Undefined(iso))
	if err == nil {
		t.Errorf("expected an error, got none")
	}
	got := *(err.(*v8.JSError))
	want := v8.JSError{Message: "error", Location: "script.js:1:21"}
	if got != want {
		t.Errorf("want %+v, got: %+v", want, got)
	}
}

func TestFunctionSourceMapUrl(t *testing.T) {
	t.Parallel()

	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()
	_, err := ctx.RunScript("function add(a, b) { return a + b; }; //# sourceMappingURL=main.js.map", "main.js")
	failIf(t, err)
	addValue, err := ctx.Global().Get("add")
	failIf(t, err)

	fn, _ := addValue.AsFunction()

	resultVal := fn.SourceMapUrl()
	if resultVal.String() != "main.js.map" {
		t.Errorf("expected main.js.map, got %v", resultVal.String())
	}

	_, err = ctx.RunScript("function sub(a, b) { return a - b; };", "")
	failIf(t, err)
	subValue, err := ctx.Global().Get("sub")
	failIf(t, err)

	subFn, _ := subValue.AsFunction()
	resultVal = subFn.SourceMapUrl()
	if !resultVal.IsUndefined() {
		t.Errorf("expected undefined, got: %v", resultVal.DetailString())
	}
}

func TestFunctionNewInstance(t *testing.T) {
	t.Parallel()

	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	iso := ctx.Isolate()

	value, err := ctx.Global().Get("Error")
	failIf(t, err)
	fn, err := value.AsFunction()
	failIf(t, err)
	messageObj, err := v8.NewValue(iso, "test message")
	failIf(t, err)
	errObj, err := fn.NewInstance(messageObj)
	failIf(t, err)

	message, err := errObj.Get("message")
	failIf(t, err)
	if !message.IsString() {
		t.Error("missing error message")
	}
	want := "test message"
	got := message.String()
	if got != want {
		t.Errorf("want %+v, got: %+v", want, got)
	}
}

func TestFunctionNewInstanceError(t *testing.T) {
	t.Parallel()

	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	_, err := ctx.RunScript("function throws() { throw 'error'; }", "script.js")
	failIf(t, err)
	throwsValue, err := ctx.Global().Get("throws")
	failIf(t, err)
	fn, _ := throwsValue.AsFunction()

	_, err = fn.NewInstance()
	if err == nil {
		t.Errorf("expected an error, got none")
	}
	got := *(err.(*v8.JSError))
	want := v8.JSError{Message: "error", Location: "script.js:1:21"}
	if got != want {
		t.Errorf("want %+v, got: %+v", want, got)
	}
}
