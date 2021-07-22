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

	"github.com/Shopify/v8go"
)

func TestFunctionTemplate(t *testing.T) {
	t.Parallel()

	if _, err := v8go.NewFunctionTemplate(nil, func(*v8go.FunctionCallbackInfo) *v8go.Value { return nil }); err == nil {
		t.Error("expected error but got <nil>")
	}

	iso, _ := v8go.NewIsolate()
	if _, err := v8go.NewFunctionTemplate(iso, nil); err == nil {
		t.Error("expected error but got <nil>")
	}

	fn, err := v8go.NewFunctionTemplate(iso, func(*v8go.FunctionCallbackInfo) *v8go.Value { return nil })
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if fn == nil {
		t.Error("expected FunctionTemplate, but got <nil>")
	}
}

func TestFunctionTemplateGetFunction(t *testing.T) {
	t.Parallel()

	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewContext(iso)

	var args *v8go.FunctionCallbackInfo
	tmpl, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args = info
		reply, _ := v8go.NewValue(iso, "hello")
		return reply
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
	printfn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		fmt.Printf("%+v\n", info.Args())
		return nil
	})
	global.Set("print", printfn, v8go.ReadOnly)
	ctx, _ := v8go.NewContext(iso, global)
	ctx.RunScript("print('foo', 'bar', 0, 1)", "")
	// Output:
	// [foo bar 0 1]
}

func ExampleFunctionTemplate_fetch() {
	iso, _ := v8go.NewIsolate()
	global, _ := v8go.NewObjectTemplate(iso)

	fetchfn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
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

	ctx, _ := v8go.NewContext(iso, global)
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
