// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"testing"
	"time"

	v8 "rogchap.com/v8go"
)

func TestIsolateTermination(t *testing.T) {
	t.Parallel()
	iso := v8.NewIsolate()
	defer iso.Dispose()

	if iso.IsExecutionTerminating() {
		t.Error("expected no execution to be terminating")
	}

	var terminating bool
	fooFn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		loop, _ := info.Args()[0].AsFunction()
		loop.Call(v8.Undefined(iso))

		terminating = iso.IsExecutionTerminating()
		return nil
	})

	global := v8.NewObjectTemplate(iso)
	global.Set("foo", fooFn)

	ctx := v8.NewContext(iso, global)
	defer ctx.Close()

	go func() {
		// [RC] find a better way to know when a script has started execution
		time.Sleep(time.Millisecond)
		iso.TerminateExecution()
	}()

	script := `function loop() { while (true) { } }; foo(loop);`
	_, e := ctx.CompileAndRun(script, "forever.js")
	if e == nil || !strings.HasPrefix(e.Error(), "ExecutionTerminated") {
		t.Errorf("unexpected error: %v", e)
	}

	if !terminating {
		t.Error("expected execution to have been terminating in function")
	}
}

func TestIsolateCompileScript(t *testing.T) {
	s := "function foo() { return 'bar'; }; foo()"

	i1 := v8.NewIsolate()
	defer i1.Dispose()

	i2 := v8.NewIsolate()
	defer i2.Dispose()
	c2 := v8.NewContext(i2)
	defer c2.Close()

	d1, err := i1.CompileScript(s, "script.js", v8.ScriptCompilerCompileOptionEagerCompile)
	fatalIf(t, err)

	val, err := c2.RunScript(s, d1, "script.js")
	fatalIf(t, err)
	if val.String() != "bar" {
		t.Fatalf("invalid value returned, expected bar got %v", val)
	}

	_, err = i1.CompileScript("invalid js", "filename", v8.ScriptCompilerCompileOptionEagerCompile)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestIsolateGetHeapStatistics(t *testing.T) {
	t.Parallel()
	iso := v8.NewIsolate()
	defer iso.Dispose()
	ctx1 := v8.NewContext(iso)
	defer ctx1.Close()
	ctx2 := v8.NewContext(iso)
	defer ctx2.Close()

	hs := iso.GetHeapStatistics()

	if hs.NumberOfNativeContexts != 3 {
		t.Error("expect NumberOfNativeContexts return 3, got", hs.NumberOfNativeContexts)
	}

	if hs.NumberOfDetachedContexts != 0 {
		t.Error("expect NumberOfDetachedContexts return 0, got", hs.NumberOfDetachedContexts)
	}
}

func TestCallbackRegistry(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	defer iso.Dispose()
	cb := func(*v8.FunctionCallbackInfo) *v8.Value { return nil }

	cb0 := iso.GetCallback(0)
	if cb0 != nil {
		t.Error("expected callback function to be <nil>")
	}
	ref1 := iso.RegisterCallback(cb)
	if ref1 != 1 {
		t.Errorf("expected callback ref == 1, got %d", ref1)
	}
	cb1 := iso.GetCallback(1)
	if fmt.Sprintf("%p", cb1) != fmt.Sprintf("%p", cb) {
		t.Errorf("unexpected callback function; want %p, got %p", cb, cb1)
	}
}

func TestIsolateDispose(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	if iso.GetHeapStatistics().TotalHeapSize == 0 {
		t.Error("Isolate incorrectly allocated")
	}

	iso.Dispose()
	// noop when called multiple times
	iso.Dispose()
	// deprecated
	iso.Close()

	if iso.GetHeapStatistics().TotalHeapSize != 0 {
		t.Error("Isolate not disposed correctly")
	}
}

func TestIsolateGarbageCollection(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	val, _ := v8.NewValue(iso, "some string")
	fmt.Println(val.String())

	tmpl := v8.NewObjectTemplate(iso)
	tmpl.Set("foo", "bar")
	v8.NewContext(iso, tmpl)

	iso.Dispose()

	runtime.GC()

	time.Sleep(time.Second)
}

func BenchmarkIsolateInitialization(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		vm := v8.NewIsolate()
		vm.Close() // force disposal of the VM
	}
}

func BenchmarkIsolateInitAndRun(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		vm := v8.NewIsolate()
		ctx := v8.NewContext(vm)
		ctx.CompileAndRun(script, "main.js")
		str, _ := json.Marshal(makeObject())
		cmd := fmt.Sprintf("process(%s)", str)
		ctx.CompileAndRun(cmd, "cmd.js")
		ctx.Close()
		vm.Close() // force disposal of the VM
	}
}

const script = `
	const process = (record) => {
		const res = [];
		for (let [k, v] of Object.entries(record)) {
			res.push({
				name: k,
				value: v,
			});
		}
		return JSON.stringify(res);
	};
`

func makeObject() interface{} {
	return map[string]interface{}{
		"a": rand.Intn(1000000),
		"b": "AAAABBBBAAAABBBBAAAABBBBAAAABBBBAAAABBBB",
	}
}
