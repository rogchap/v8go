// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"rogchap.com/v8go"
)

func TestFunctionTemplate(t *testing.T) {
	t.Parallel()

	if _, err := v8go.NewFunctionTemplate(nil, func(*v8go.FunctionCallbackInfo) (v8go.Valuer, error) { return nil, nil }); err == nil {
		t.Error("expected error but got <nil>")
	}

	iso, _ := v8go.NewIsolate()
	if _, err := v8go.NewFunctionTemplate(iso, nil); err == nil {
		t.Error("expected error but got <nil>")
	}

	fn, err := v8go.NewFunctionTemplate(iso, func(*v8go.FunctionCallbackInfo) (v8go.Valuer, error) { return nil, nil })
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if fn == nil {
		t.Error("expected FunctionTemplate, but got <nil>")
	}
}

func TestFunctionTemplateError(t *testing.T) {
	t.Parallel()

	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewExecContext(iso)

	tmpl, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) (v8go.Valuer, error) {
		return nil, errors.New("test")
	})
	fn := tmpl.GetFunction(ctx)
	_, err := fn.Call()
	if err == nil {
		t.Error("function should throw new exception")
	}
}

func TestFunctionTemplateGetFunction(t *testing.T) {
	t.Parallel()

	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewExecContext(iso)

	var args *v8go.FunctionCallbackInfo
	tmpl, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) (v8go.Valuer, error) {
		args = info
		reply, _ := v8go.NewValue(iso, "hello")
		return reply, nil
	})
	fn := tmpl.GetFunction(ctx)
	ten, err := v8go.NewValue(iso, int32(10))
	if err != nil {
		t.Fatal(err)
	}
	ret, err := fn.Call(ten)
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

func ExampleFunctionTemplate() {
	iso, _ := v8go.NewIsolate()
	global, _ := v8go.NewObjectTemplate(iso)
	printfn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) (v8go.Valuer, error) {
		fmt.Printf("%+v\n", info.Args())
		return nil, nil
	})
	global.Set("print", printfn, v8go.ReadOnly)
	ctx, _ := v8go.NewExecContext(iso, global)
	ctx.RunScript("print('foo', 'bar', 0, 1)", "")
	// Output:
	// [foo bar 0 1]
}

func ExampleFunctionTemplate_promise() {
	iso, _ := v8go.NewIsolate()
	global, _ := v8go.NewObjectTemplate(iso)

	fn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) (v8go.Valuer, error) {
		resolver, _ := v8go.NewPromiseResolver(info.ExecContext())

		go func() {
			val, _ := v8go.NewValue(iso, "ZOMGBBQ it works!")
			resolver.Resolve(val)
		}()
		return resolver.GetPromise().Value, nil
	})
	global.Set("resolve", fn, v8go.ReadOnly)

	ctx, _ := v8go.NewExecContext(iso, global)
	val, _ := ctx.RunScript("resolve()", "")
	prom, _ := val.AsPromise()

	// wait for the promise to resolve
	for prom.State() == v8go.Pending {
		continue
	}
	fmt.Printf("%s\n", strings.Split(prom.Result().String(), "\n")[0])
	// Output:
	// ZOMGBBQ it works!
}
