// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"rogchap.com/v8go"
)

func TestContextRunScript(t *testing.T) {
	t.Parallel()
	ctx, _ := v8go.NewContext(nil)
	ctx.RunScript(`const add = (a, b) => a + b`, "add.js")
	val, _ := ctx.RunScript(`add(3, 4)`, "main.js")
	rtn := val.String()
	if rtn != "7" {
		t.Errorf("script returned an unexpected value: expected %q, got %q", "7", rtn)
	}

	_, err := ctx.RunScript(`add`, "func.js")
	if err != nil {
		t.Errorf("error not expected: %v", err)
	}

	iso, _ := ctx.Isolate()
	ctx2, _ := v8go.NewContext(iso)
	_, err = ctx2.RunScript(`add`, "ctx2.js")
	if err == nil {
		t.Error("error expected but was <nil>")
	}
}

func TestRunScriptJSExceptions(t *testing.T) {
	t.Parallel()

	tests := [...]struct {
		name   string
		source string
		origin string
		err    string
	}{
		{"SyntaxError", "bad js syntax", "syntax.js", "SyntaxError: Unexpected identifier"},
		{"ReferenceError", "add()", "add.js", "ReferenceError: add is not defined"},
	}

	ctx, _ := v8go.NewContext(nil)
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := ctx.RunScript(tt.source, tt.origin)
			if err == nil {
				t.Error("error expected but got <nil>")
				return
			}
			if err.Error() != tt.err {
				t.Errorf("expected %q, got %q", tt.err, err.Error())
			}
		})
	}
}

func TestContextRunModule(t *testing.T) {
	t.Parallel()
	ctx, _ := v8go.NewContext(nil)
	val, err := ctx.RunModule(`export function add(a, b) { return a + b }`, "add.mjs")
	if err != nil {
		t.Errorf("RunModule returned an error: %+v", err)
		return
	}
	if !val.IsModuleNamespaceObject() {
		t.Errorf("RunModule returned an unexpected value: %+v", val.DetailString())
		return
	}
	obj, _ := val.AsObject()
	add, _ := obj.Get("add")
	if !add.IsFunction() {
		t.Errorf("expected to get exported 'add' function, got: %+v", add.DetailString())
		return
	}
	fn, _ := add.AsFunction()
	iso, _ := ctx.Isolate()
	arg1, _ := v8go.NewValue(iso, int32(1))
	resultValue, _ := fn.Call(arg1, arg1)
	if resultValue.Int32() != 2 {
		t.Errorf("expected 1 + 1 = 2, got: %v", resultValue.DetailString())
	}
}

func TestRunModuleJSExceptions(t *testing.T) {
	t.Parallel()

	tests := [...]struct {
		name   string
		source string
		origin string
		err    string
	}{
		{"SyntaxError", "bad js syntax", "syntax.js", "SyntaxError: Unexpected identifier"},
		{"ReferenceError", "add()", "add.js", "ReferenceError: add is not defined"},
		{"import", "import { add } from 'dep.js'", "import.js", "import not supported"},
	}

	ctx, _ := v8go.NewContext(nil)
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := ctx.RunModule(tt.source, tt.origin)
			if err == nil {
				t.Error("error expected but got <nil>")
				return
			}
			if err.Error() != tt.err {
				t.Errorf("expected %q, got %q", tt.err, err.Error())
			}
		})
	}
}


// TestRunModuleTopLevelAwait verifies that top-level await works in modules.
// It exercises a different module resolution control flow.
func TestRunModuleTopLevelAwait(t *testing.T) {
	v8go.SetFlags("--harmony_top_level_await")
	defer v8go.SetFlags("--noharmony_top_level_await")

	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewContext(iso)

	t.Run("success", func(t *testing.T) {
		val, err := ctx.RunModule(`
function resolvePromise() {
  return new Promise(resolve => {
    resolve('resolved');
  });
}

export async function asyncCall() {
  return await resolvePromise();
}

await asyncCall();
`, "")
		if err != nil {
			t.Errorf("%+v", err)
		}
		obj, _ := val.AsObject()
		asyncCall, _ := obj.Get("asyncCall")
		if !asyncCall.IsAsyncFunction() {
			t.Errorf("expected async function, was: %v", asyncCall.DetailString())
		}
	})

	t.Run("error", func(t *testing.T) {
		_, err := ctx.RunModule(`
function resolvePromise() {
  return new Promise((resolve, reject) => {
    reject('rejected');
  });
}

export async function asyncCall() {
  return await resolvePromise();
}

await asyncCall();
`, "")
		if err == nil {
			t.Errorf("expected error, got none")
		}
		want := "rejected"
		if got := err.Error(); got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})
}


func TestContextRegistry(t *testing.T) {
	t.Parallel()

	ctx, _ := v8go.NewContext()
	ctxref := ctx.Ref()

	c1 := v8go.GetContext(ctxref)
	if c1 != nil {
		t.Error("expected context to be <nil>")
	}

	ctx.Register()
	c2 := v8go.GetContext(ctxref)
	if c2 == nil {
		t.Error("expected context, but got <nil>")
	}
	if c2 != ctx {
		t.Errorf("contexts should match %p != %p", c2, ctx)
	}
	ctx.Deregister()

	c3 := v8go.GetContext(ctxref)
	if c3 != nil {
		t.Error("expected context to be <nil>")
	}
}

func TestMemoryLeak(t *testing.T) {
	t.Parallel()

	iso, _ := v8go.NewIsolate()

	for i := 0; i < 6000; i++ {
		ctx, _ := v8go.NewContext(iso)
		obj := ctx.Global()
		_ = obj.String()
		_, _ = ctx.RunScript("2", "")
		ctx.Close()
	}
	if n := iso.GetHeapStatistics().NumberOfNativeContexts; n >= 6000 {
		t.Errorf("Context not being GC'd, got %d native contexts", n)
	}
}

func BenchmarkContext(b *testing.B) {
	b.ReportAllocs()
	vm, _ := v8go.NewIsolate()
	defer vm.Close()
	for n := 0; n < b.N; n++ {
		ctx, _ := v8go.NewContext(vm)
		ctx.RunScript(script, "main.js")
		str, _ := json.Marshal(makeObject())
		cmd := fmt.Sprintf("process(%s)", str)
		ctx.RunScript(cmd, "cmd.js")
		ctx.Close()
	}
}

func ExampleContext() {
	ctx, _ := v8go.NewContext()
	ctx.RunScript("const add = (a, b) => a + b", "math.js")
	ctx.RunScript("const result = add(3, 4)", "main.js")
	val, _ := ctx.RunScript("result", "value.js")
	fmt.Println(val)
	// Output:
	// 7
}

func ExampleContext_isolate() {
	iso, _ := v8go.NewIsolate()
	ctx1, _ := v8go.NewContext(iso)
	ctx1.RunScript("const foo = 'bar'", "context_one.js")
	val, _ := ctx1.RunScript("foo", "foo.js")
	fmt.Println(val)

	ctx2, _ := v8go.NewContext(iso)
	_, err := ctx2.RunScript("foo", "context_two.js")
	fmt.Println(err)
	// Output:
	// bar
	// ReferenceError: foo is not defined
}

func ExampleContext_globalTemplate() {
	iso, _ := v8go.NewIsolate()
	obj, _ := v8go.NewObjectTemplate(iso)
	obj.Set("version", "v1.0.0")
	ctx, _ := v8go.NewContext(iso, obj)
	val, _ := ctx.RunScript("version", "main.js")
	fmt.Println(val)
	// Output:
	// v1.0.0
}
