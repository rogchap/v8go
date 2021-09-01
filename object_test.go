// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"fmt"
	"testing"

	"rogchap.com/v8go"
)

func TestObjectMethodCall(t *testing.T) {
	t.Parallel()

	ctx, err := v8go.NewContext()
	failIf(t, err)
	val, err := ctx.RunScript(`class Obj { constructor(input) { this.input = input } print() { return this.input.toString() } }; new Obj("some val")`, "")
	failIf(t, err)
	obj, err := val.AsObject()
	failIf(t, err)
	val, err = obj.MethodCall("print")
	failIf(t, err)
	if val.String() != "some val" {
		t.Errorf("unexpected value: %q", val)
	}
	_, err = obj.MethodCall("nope")
	if err == nil {
		t.Errorf("expected an error, got none")
	}
}

func TestObjectSet(t *testing.T) {
	t.Parallel()

	ctx, _ := v8go.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()
	val, _ := ctx.RunScript("const foo = {}; foo", "")
	obj, _ := val.AsObject()
	obj.Set("bar", "baz")
	baz, _ := ctx.RunScript("foo.bar", "")
	if baz.String() != "baz" {
		t.Errorf("unexpected value: %q", baz)
	}
	if err := obj.Set("", nil); err == nil {
		t.Error("expected error but got <nil>")
	}
	if err := obj.Set("a", 0); err == nil {
		t.Error("expected error but got <nil>")
	}
	obj.SetIdx(10, "ten")
	if ten, _ := ctx.RunScript("foo[10]", ""); ten.String() != "ten" {
		t.Errorf("unexpected value: %q", ten)
	}
}

func TestObjectGet(t *testing.T) {
	t.Parallel()

	ctx, _ := v8go.NewContext()
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

	ctx, _ := v8go.NewContext()
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

	ctx, _ := v8go.NewContext()
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
	iso, _ := v8go.NewIsolate()
	defer iso.Dispose()
	ctx, _ := v8go.NewContext(iso)
	defer ctx.Close()
	global := ctx.Global()

	console := v8go.NewObjectTemplate(iso)
	logfn := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
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
