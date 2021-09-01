// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"rogchap.com/v8go"
)

func TestFunctionTemplate(t *testing.T) {
	t.Parallel()

	iso := v8go.NewIsolate()
	defer iso.Dispose()
	fn := v8go.NewFunctionTemplate(iso, func(*v8go.FunctionCallbackInfo) *v8go.Value { return nil })
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
	v8go.NewFunctionTemplate(nil, func(*v8go.FunctionCallbackInfo) *v8go.Value {
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
	iso := v8go.NewIsolate()
	defer iso.Dispose()
	v8go.NewFunctionTemplate(iso, nil)
}

func TestFunctionTemplateGetFunction(t *testing.T) {
	t.Parallel()

	iso := v8go.NewIsolate()
	defer iso.Dispose()
	ctx := v8go.NewContext(iso)
	defer ctx.Close()

	var args *v8go.FunctionCallbackInfo
	tmpl := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args = info
		reply, _ := v8go.NewValue(iso, "hello")
		return reply
	})
	fn := tmpl.GetFunction(ctx)
	ten, err := v8go.NewValue(iso, int32(10))
	if err != nil {
		t.Fatal(err)
	}
	ret, err := fn.Call(ctx.Global(), ten)
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

	iso := v8go.NewIsolate()

	foo := v8go.NewObjectTemplate(iso)
	foo.Set("name", "foobar")

	var this *v8go.Object
	barfn := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		this = info.This()
		return nil
	})
	foo.Set("bar", barfn)

	global := v8go.NewObjectTemplate(iso)
	global.Set("foo", foo)

	ctx := v8go.NewContext(iso, global)
	ctx.RunScript("foo.bar()", "")

	v, _ := this.Get("name")
	if v.String() != "foobar" {
		t.Errorf("expected this.name to be foobar, but got %q", v)
	}
}

func ExampleFunctionTemplate() {
	iso := v8go.NewIsolate()
	defer iso.Dispose()
	global := v8go.NewObjectTemplate(iso)
	printfn := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		fmt.Printf("%+v\n", info.Args())
		return nil
	})
	global.Set("print", printfn, v8go.ReadOnly)
	ctx := v8go.NewContext(iso, global)
	defer ctx.Close()
	ctx.RunScript("print('foo', 'bar', 0, 1)", "")
	// Output:
	// [foo bar 0 1]
}

func ExampleFunctionTemplate_fetch() {
	iso := v8go.NewIsolate()
	defer iso.Dispose()
	global := v8go.NewObjectTemplate(iso)

	fetchfn := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		url := args[0].String()

		resolver, _ := v8go.NewPromiseResolver(info.Context())

		go func() {
			res, _ := http.Get(url)
			body, _ := ioutil.ReadAll(res.Body)
			val, _ := v8go.NewValue(iso, string(body))
			resolver.Resolve(val)
		}()
		return resolver.GetPromise().Value
	})
	global.Set("fetch", fetchfn, v8go.ReadOnly)

	ctx := v8go.NewContext(iso, global)
	defer ctx.Close()
	val, _ := ctx.RunScript("fetch('https://rogchap.com/v8go')", "")
	prom, _ := val.AsPromise()

	// wait for the promise to resolve
	for prom.State() == v8go.Pending {
		continue
	}
	fmt.Printf("%s\n", strings.Split(prom.Result().String(), "\n")[0])
	// Output:
	// <!DOCTYPE html>
}
