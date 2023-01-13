// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"fmt"
	"testing"

	v8 "rogchap.com/v8go"
)

func TestFunctionTemplate(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	defer iso.Dispose()
	fn := v8.NewFunctionTemplate(iso, func(*v8.FunctionCallbackInfo) *v8.Value { return nil })
	if fn == nil {
		t.Error("expected FunctionTemplate, but got <nil>")
	}
}

func TestFunctionTemplate_panic_on_nil_isolate(t *testing.T) {
	t.Parallel()

	defer func() {
		if err := recover(); err == nil {
			t.Error("expected panic")
		}
	}()
	v8.NewFunctionTemplate(nil, func(*v8.FunctionCallbackInfo) *v8.Value {
		t.Error("unexpected call")
		return nil
	})
}

func TestFunctionTemplate_panic_on_nil_callback(t *testing.T) {
	t.Parallel()

	defer func() {
		if err := recover(); err == nil {
			t.Error("expected panic")
		}
	}()
	iso := v8.NewIsolate()
	defer iso.Dispose()
	v8.NewFunctionTemplate(iso, nil)
}
func TestFunctionTemplate_generates_values(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	defer iso.Dispose()
	global := v8.NewObjectTemplate(iso)
	printfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		fmt.Printf("%+v\n", info.Args())
		return nil
	})
	global.Set("print", printfn, v8.ReadOnly)
	ctx := v8.NewContext(iso, global)
	defer ctx.Close()
	ctx.RunScript("print('foo', 'bar', 0, 1)", "")
	if ctx.RetainedValueCount() != 5 {
		t.Errorf("expected 5 retained values, got: %d", ctx.RetainedValueCount())
	}
}

func TestFunctionTemplate_releases_values(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	defer iso.Dispose()
	global := v8.NewObjectTemplate(iso)
	printfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		defer info.Release()
		fmt.Printf("%+v\n", info.Args())
		return nil
	})
	global.Set("print", printfn, v8.ReadOnly)
	ctx := v8.NewContext(iso, global)
	defer ctx.Close()
	ctx.RunScript("print('foo', 'bar', 0, 1)", "")
	if ctx.RetainedValueCount() != 0 {
		t.Errorf("expected 0 retained values, got: %d", ctx.RetainedValueCount())
	}
}

func TestFunctionTemplateGetFunction(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	defer iso.Dispose()
	ctx := v8.NewContext(iso)
	defer ctx.Close()

	var args *v8.FunctionCallbackInfo
	tmpl := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		args = info
		reply, _ := v8.NewValue(iso, "hello")
		return reply
	})
	fn := tmpl.GetFunction(ctx)
	ten, err := v8.NewValue(iso, int32(10))
	if err != nil {
		t.Fatal(err)
	}
	ret, err := fn.Call(v8.Undefined(iso), ten)
	if err != nil {
		t.Fatal(err)
	}
	if len(args.Args()) != 1 || args.Args()[0].String() != "10" {
		t.Fatalf("expected args [10], got: %+v", args.Args())
	}
	if !ret.IsString() || ret.String() != "hello" {
		t.Fatalf("expected return value of 'hello', was: %v", ret)
	}
}

func TestFunctionCallbackInfoThis(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	defer iso.Dispose()

	foo := v8.NewObjectTemplate(iso)
	foo.Set("name", "foobar")

	var this *v8.Object
	barfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		this = info.This()
		return nil
	})
	foo.Set("bar", barfn)

	global := v8.NewObjectTemplate(iso)
	global.Set("foo", foo)

	ctx := v8.NewContext(iso, global)
	defer ctx.Close()
	ctx.RunScript("foo.bar()", "")

	v, _ := this.Get("name")
	if v.String() != "foobar" {
		t.Errorf("expected this.name to be foobar, but got %q", v)
	}
}

func ExampleFunctionTemplate() {
	iso := v8.NewIsolate()
	defer iso.Dispose()
	global := v8.NewObjectTemplate(iso)
	printfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		fmt.Printf("%+v\n", info.Args())
		return nil
	})
	global.Set("print", printfn, v8.ReadOnly)
	ctx := v8.NewContext(iso, global)
	defer ctx.Close()
	ctx.RunScript("print('foo', 'bar', 0, 1)", "")
	// Output:
	// [foo bar 0 1]
}
