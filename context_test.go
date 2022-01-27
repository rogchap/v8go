// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"encoding/json"
	"fmt"
	"testing"

	v8 "rogchap.com/v8go"
)

func TestContextExec(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext(nil)
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

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

	iso := ctx.Isolate()
	ctx2 := v8.NewContext(iso)
	defer ctx2.Close()
	_, err = ctx2.RunScript(`add`, "ctx2.js")
	if err == nil {
		t.Error("error expected but was <nil>")
	}
}

func TestNewContextFromSnapshotErrorWhenIsolateHasNoStartupData(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	defer iso.Dispose()

	ctx, err := v8.NewContextFromSnapshot(iso, 1)

	if ctx != nil {
		t.Errorf("expected nil context got: %+v", ctx)
	}
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

	ctx := v8.NewContext(nil)
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	ctxref := ctx.Ref()

	c1 := v8.GetContext(ctxref)
	if c1 == nil {
		t.Error("expected context, but got <nil>")
	}
	if c1 != ctx {
		t.Errorf("contexts should match %p != %p", c1, ctx)
	}

	ctx.Close()

	c2 := v8.GetContext(ctxref)
	if c2 != nil {
		t.Error("expected context to be <nil> after close")
	}
}

func TestMemoryLeak(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	defer iso.Dispose()

	for i := 0; i < 6000; i++ {
		ctx := v8.NewContext(iso)
		obj := ctx.Global()
		_ = obj.String()
		_, _ = ctx.RunScript("2", "")
		ctx.Close()
	}
	if n := iso.GetHeapStatistics().NumberOfNativeContexts; n >= 6000 {
		t.Errorf("Context not being GC'd, got %d native contexts", n)
	}
}

// https://github.com/rogchap/v8go/issues/186
func TestRegistryFromJSON(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	defer iso.Dispose()

	global := v8.NewObjectTemplate(iso)
	err := global.Set("location", v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		v, err := v8.NewValue(iso, "world")
		fatalIf(t, err)
		return v
	}))
	fatalIf(t, err)

	ctx := v8.NewContext(iso, global)
	defer ctx.Close()

	v, err := ctx.RunScript(`
		new Proxy({
			"hello": "unknown"
		}, {
			get: function () {
				return location()
			},
		})
	`, "main.js")
	fatalIf(t, err)

	s, err := v8.JSONStringify(ctx, v)
	fatalIf(t, err)

	expected := `{"hello":"world"}`
	if s != expected {
		t.Fatalf("expected %q, got %q", expected, s)
	}
}

func BenchmarkContext(b *testing.B) {
	b.ReportAllocs()
	iso := v8.NewIsolate()
	defer iso.Dispose()
	for n := 0; n < b.N; n++ {
		ctx := v8.NewContext(iso)
		ctx.RunScript(script, "main.js")
		str, _ := json.Marshal(makeObject())
		cmd := fmt.Sprintf("process(%s)", str)
		ctx.RunScript(cmd, "cmd.js")
		ctx.Close()
	}
}

func ExampleContext() {
	ctx := v8.NewContext()
	defer ctx.Isolate().Dispose()
	defer ctx.Close()
	ctx.RunScript("const add = (a, b) => a + b", "math.js")
	ctx.RunScript("const result = add(3, 4)", "main.js")
	val, _ := ctx.RunScript("result", "value.js")
	fmt.Println(val)
	// Output:
	// 7
}

func ExampleContext_isolate() {
	iso := v8.NewIsolate()
	defer iso.Dispose()
	ctx1 := v8.NewContext(iso)
	defer ctx1.Close()
	ctx1.RunScript("const foo = 'bar'", "context_one.js")
	val, _ := ctx1.RunScript("foo", "foo.js")
	fmt.Println(val)

	ctx2 := v8.NewContext(iso)
	defer ctx2.Close()
	_, err := ctx2.RunScript("foo", "context_two.js")
	fmt.Println(err)
	// Output:
	// bar
	// ReferenceError: foo is not defined
}

func ExampleContext_globalTemplate() {
	iso := v8.NewIsolate()
	defer iso.Dispose()
	obj := v8.NewObjectTemplate(iso)
	obj.Set("version", "v1.0.0")
	ctx := v8.NewContext(iso, obj)
	defer ctx.Close()
	val, _ := ctx.RunScript("version", "main.js")
	fmt.Println(val)
	// Output:
	// v1.0.0
}
