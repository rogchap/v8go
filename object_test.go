// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"errors"
	"fmt"
	"testing"

	"rogchap.com/v8go"
	v8 "rogchap.com/v8go"
)

func TestObjectMethodCall(t *testing.T) {
	t.Parallel()

	ctx := v8.NewContext()
	iso := ctx.Isolate()
	defer iso.Dispose()
	defer ctx.Close()
	val, _ := ctx.RunScript(`class Obj { constructor(input) { this.input = input, this.prop = "" } print() { return this.input.toString() } }; new Obj("some val")`, "")
	obj, _ := val.AsObject()
	val, err := obj.MethodCall("print")
	fatalIf(t, err)
	if val.String() != "some val" {
		t.Errorf("unexpected value: %q", val)
	}
	_, err = obj.MethodCall("prop")
	if err == nil {
		t.Errorf("expected an error, got none")
	}

	val, err = ctx.RunScript(`class Obj2 { print(str) { return str.toString() }; get fails() { throw "error" } }; new Obj2()`, "")
	fatalIf(t, err)
	obj, _ = val.AsObject()
	arg, _ := v8.NewValue(iso, "arg")
	val, err = obj.MethodCall("print", arg)
	fatalIf(t, err)
	if val.String() != "arg" {
		t.Errorf("unexpected value: %q", val)
	}
	_, err = obj.MethodCall("fails")
	if err == nil {
		t.Errorf("expected an error, got none")
	}
}

func TestObjectSet(t *testing.T) {
	t.Parallel()

	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()
	val, _ := ctx.RunScript("const foo = {}; foo", "")
	obj, _ := val.AsObject()
	if err := obj.Set("bar", "baz"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	baz, _ := ctx.RunScript("foo.bar", "")
	if baz.String() != "baz" {
		t.Errorf("unexpected value: %q", baz)
	}

	if err := obj.Set("", "zero"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	val, err := ctx.RunScript("foo['']", "")
	fatalIf(t, err)
	if val.String() != "zero" {
		t.Errorf("unexpected value: %q", val)
	}

	if err := obj.Set("a", nil); err == nil {
		t.Error("expected error but got <nil>")
	}
	if err := obj.Set("a", 0); err == nil {
		t.Error("expected error but got <nil>")
	}
	if err := obj.SetIdx(10, "ten"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := obj.SetIdx(10, t); err == nil {
		t.Error("expected error but got <nil>")
	}
	if ten, _ := ctx.RunScript("foo[10]", ""); ten.String() != "ten" {
		t.Errorf("unexpected value: %q", ten)
	}
}

func TestObjectInternalFields(t *testing.T) {
	iso := v8.NewIsolate()
	defer iso.Dispose()
	ctx := v8.NewContext(iso)
	defer ctx.Close()

	tmpl := v8.NewObjectTemplate(iso)
	obj, err := tmpl.NewInstance(ctx)
	fatalIf(t, err)
	if count := obj.InternalFieldCount(); count != 0 {
		t.Errorf("expected 0 got %v", count)
	}
	if recoverPanic(func() { obj.GetInternalField(0) }) == nil {
		t.Error("expected panic")
	}

	tmpl = v8.NewObjectTemplate(iso)
	tmpl.SetInternalFieldCount(1)
	if count := tmpl.InternalFieldCount(); count != 1 {
		t.Errorf("expected 1 got %v", count)
	}

	obj, err = tmpl.NewInstance(ctx)
	fatalIf(t, err)
	if count := obj.InternalFieldCount(); count != 1 {
		t.Errorf("expected 1 got %v", count)
	}

	if v := obj.GetInternalField(0); !v.SameValue(v8.Undefined(iso)) {
		t.Errorf("unexpected value: %q", v)
	}

	if err := obj.SetInternalField(0, t); err == nil {
		t.Error("expected unsupported value error")
	}

	err = obj.SetInternalField(0, "baz")
	fatalIf(t, err)
	if v := obj.GetInternalField(0); v.String() != "baz" {
		t.Errorf("unexpected value: %q", v)
	}

	if recoverPanic(func() { obj.SetInternalField(1, "baz") }) == nil {
		t.Error("expected panic from index out of bounds")
	}
}

func TestObjectGet(t *testing.T) {
	t.Parallel()

	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()
	val, _ := ctx.RunScript("const foo = { bar: 'baz'}; foo", "")
	obj, _ := val.AsObject()
	if bar, _ := obj.Get("bar"); bar.String() != "baz" {
		t.Errorf("unexpected value: %q", bar)
	}
	if baz, _ := obj.Get("baz"); !baz.IsUndefined() {
		t.Errorf("unexpected value: %q", baz)
	}
	ctx.RunScript("foo[5] = 5", "")
	if five, _ := obj.GetIdx(5); five.Integer() != 5 {
		t.Errorf("unexpected value: %q", five)
	}
	if u, _ := obj.GetIdx(55); !u.IsUndefined() {
		t.Errorf("unexpected value: %q", u)
	}
}

func TestObjectHas(t *testing.T) {
	t.Parallel()

	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()
	val, _ := ctx.RunScript("const foo = {a: 1, '2': 2}; foo", "")
	obj, _ := val.AsObject()
	if !obj.Has("a") {
		t.Error("expected true, got false")
	}
	if obj.Has("c") {
		t.Error("expected false, got true")
	}
	if !obj.HasIdx(2) {
		t.Error("expected true, got false")
	}
	if obj.HasIdx(1) {
		t.Error("expected false, got true")
	}
}

func TestObjectDelete(t *testing.T) {
	t.Parallel()

	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()
	val, _ := ctx.RunScript("const foo = { bar: 'baz', '2': 2}; foo", "")
	obj, _ := val.AsObject()
	if !obj.Has("bar") {
		t.Error("expected property to exist")
	}
	if !obj.Delete("bar") {
		t.Error("expected delete to return true, got false")
	}
	if obj.Has("bar") {
		t.Error("expected property to be deleted")
	}
	if !obj.DeleteIdx(2) {
		t.Error("expected delete to return true, got false")
	}

}

func ExampleObject_global() {
	iso := v8.NewIsolate()
	defer iso.Dispose()
	ctx := v8.NewContext(iso)
	defer ctx.Close()
	global := ctx.Global()

	console := v8.NewObjectTemplate(iso)
	logfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		fmt.Println(info.Args()[0])
		return nil
	})
	console.Set("log", logfn)
	consoleObj, _ := console.NewInstance(ctx)

	global.Set("console", consoleObj)
	ctx.RunScript("console.log('foo')", "")
	// Output:
	// foo
}

func createObjectFunctionCallback(info *v8go.FunctionCallbackInfo) *v8go.Value {
	iso := info.Context().Isolate()
	args := info.Args()
	if len(args) != 2 {
		e, _ := v8.NewValue(iso, "Function createObject expects 2 parameters")
		return iso.ThrowException(e)
	}
	if !args[0].IsInt32() || !args[1].IsInt32() {
		e, _ := v8.NewValue(iso, "Function createObject expects 2 Int32 parameters")
		return iso.ThrowException(e)
	}
	read := args[0].Int32()
	written := args[1].Int32()
	obj := v8go.NewObject(info.Context()) // create object
	obj.Set("read", read)                 // set some properties
	obj.Set("written", written)
	return obj.Value
}

func injectObjectTester(ctx *v8go.Context, funcName string, funcCb v8go.FunctionCallback) error {
	if ctx == nil {
		return errors.New("ctx is required")
	}

	iso := ctx.Isolate()

	con := v8go.NewObjectTemplate(iso)

	funcTempl := v8go.NewFunctionTemplate(iso, funcCb)

	if err := con.Set(funcName, funcTempl, v8go.ReadOnly); err != nil {
		return fmt.Errorf("ObjectTemplate.Set: %v", err)
	}

	nativeObj, err := con.NewInstance(ctx)
	if err != nil {
		return fmt.Errorf("ObjectTemplate.NewInstance: %v", err)
	}

	global := ctx.Global()

	if err := global.Set("native", nativeObj); err != nil {
		return fmt.Errorf("global.Set: %v", err)
	}

	return nil
}

// Test that golang can create an object with "read", "written" int32 properties and pass that back to JS.
func TestObjectCreate(t *testing.T) {
	t.Parallel()
	iso := v8go.NewIsolate()
	ctx := v8go.NewContext(iso)

	if err := injectObjectTester(ctx, "createObject", createObjectFunctionCallback); err != nil {
		t.Error(err)
	}

	js := `
		obj = native.createObject(123, 456);
		obj.read + obj.written;
	`

	val, err := ctx.RunScript(js, "")
	if err != nil {
		t.Errorf("Got error from script: %v", err)
	}
	if val == nil {
		t.Errorf("Got nil value from script")
	}
	if !val.IsInt32() {
		t.Errorf("Expected int32 value from script")
	}
	fmt.Printf("Script return value: %d\n", val.Int32())
	if val.Int32() != 123+456 {
		t.Errorf("Got wrong return value from script: %d", val.Int32())
	}
}
