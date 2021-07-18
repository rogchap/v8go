// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"errors"
	"fmt"
	"log"
	"testing"

	"rogchap.com/v8go"
)

func TestObjectSet(t *testing.T) {
	t.Parallel()

	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewExecContext(iso)
	val, _ := ctx.RunScript("const foo = {}; foo", "")
	obj, _ := val.AsObject()
	bazv, _ := v8go.NewValue(iso, "baz")
	obj.Set("bar", bazv)
	baz, _ := ctx.RunScript("foo.bar", "")
	if baz.String() != "baz" {
		t.Errorf("unexpected value: %q", baz)
	}
	tenv, _ := v8go.NewValue(iso, "ten")
	obj.SetIdx(10, tenv)
	if ten, _ := ctx.RunScript("foo[10]", ""); ten.String() != "ten" {
		t.Errorf("unexpected value: %q", ten)
	}
}

func TestObjectGet(t *testing.T) {
	t.Parallel()

	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewExecContext(iso)
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

	ctx, _ := v8go.NewExecContext()
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

	ctx, _ := v8go.NewExecContext()
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
	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewExecContext(iso)
	global, _ := ctx.Global()

	console, _ := v8go.NewObjectTemplate(iso)
	logfn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) (v8go.Valuer, error) {
		fmt.Println(info.Args()[0])
		return nil, nil
	})
	console.Set("log", logfn)
	consoleObj, _ := console.GetObject(ctx)

	global.Set("console", consoleObj)
	ctx.RunScript("console.log('foo')", "")
	// Output:
	// foo
}

func createObjectFunctionCallback(info *v8go.FunctionCallbackInfo) (v8go.Valuer, error) {
	iso, err := info.ExecContext().Isolate()
	if err != nil {
		log.Fatalf("Could not get isolate from context: %v\n", err)
	}
	args := info.Args()
	if len(args) != 2 {
		iso.ThrowException("Function createObject expects 2 parameters")
		return nil, nil
	}
	if !args[0].IsInt32() || !args[1].IsInt32() {
		iso.ThrowException("Function createObject expects 2 Int32 parameters")
		return nil, nil
	}
	ctx, _ := v8go.NewExecContext(iso)

	obj := v8go.NewObject(ctx) // create object
	obj.Set("read", args[0])   // set some properties
	obj.Set("written", args[1])
	return obj.Value, nil
}

func injectObjectTester(ctx *v8go.ExecContext, funcName string, funcCb v8go.FunctionCallback) error {
	if ctx == nil {
		return errors.New("ctx is required")
	}

	iso, err := ctx.Isolate()
	if err != nil {
		return fmt.Errorf("ctx.Isolate: %v", err)
	}

	con, err := v8go.NewObjectTemplate(iso)
	if err != nil {
		return fmt.Errorf("NewObjectTemplate: %v", err)
	}

	funcTempl, err := v8go.NewFunctionTemplate(iso, funcCb)
	if err != nil {
		return fmt.Errorf("NewFunctionTemplate: %v", err)
	}

	if err := con.Set(funcName, funcTempl, v8go.ReadOnly); err != nil {
		return fmt.Errorf("ObjectTemplate.Set: %v", err)
	}

	nativeObj, err := con.GetObject(ctx)
	if err != nil {
		return fmt.Errorf("ObjectTemplate.NewInstance: %v", err)
	}

	global, _ := ctx.Global()

	if err := global.Set("native", nativeObj); err != nil {
		return fmt.Errorf("global.Set: %v", err)
	}

	return nil
}

// Test that golang can create an object with "read", "written" int32 properties and pass that back to JS.
func TestObjectCreate(t *testing.T) {
	t.Parallel()
	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewExecContext(iso)

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

func TestNewObject(t *testing.T) {
	t.Parallel()
	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewExecContext()
	ok, _ := v8go.NewValue(iso, "ok")

	obj := v8go.NewObject(ctx)
	err := obj.Set("test", ok)
	if err != nil {
		t.Errorf("Got error from setting object property: %v", err)
	}
}

func TestNewObjectWithFunctionalTemplate(t *testing.T) {
	t.Parallel()
	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewExecContext(iso)
	ok, _ := v8go.NewValue(iso, "ok")

	fn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) (v8go.Valuer, error) {
		obj := v8go.NewObject(ctx)
		err := obj.Set("test", ok)
		if err != nil {
			t.Errorf("Got error from setting object property: %v", err)
		}

		return obj.Value, nil
	})

	res, err := fn.GetFunction(ctx).Call()
	if err != nil {
		t.Errorf("Got error from calling function: %v", err)
	}

	obj, _ := res.AsObject()
	test, _ := obj.Get("test")
	if test.String() != "ok" {
		t.Errorf("functional template: new object: test needs to pass %q but returned %q", "ok", test.String())
	}
}
