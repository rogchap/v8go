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

func TestContextExec(t *testing.T) {
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

func TestJSExceptions(t *testing.T) {
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
