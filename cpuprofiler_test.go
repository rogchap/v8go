// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	"rogchap.com/v8go"
)

func TestCPUProfilerDispose(t *testing.T) {
	t.Parallel()

	iso := v8go.NewIsolate()
	defer iso.Dispose()
	cpuProfiler := v8go.NewCPUProfiler(iso)

	cpuProfiler.Dispose()
	// noop when called multiple times
	cpuProfiler.Dispose()

	// verify does not panic once disposed
	cpuProfiler.StartProfiling("")
	cpuProfiler.StopProfiling("")

	cpuProfiler = v8go.NewCPUProfiler(iso)
	defer cpuProfiler.Dispose()
	iso.Dispose()
	// verify does not panic once isolate disposed
	cpuProfiler.StartProfiling("")
	cpuProfiler.StopProfiling("")
}

func TestCPUProfiler(t *testing.T) {
	t.Parallel()

	ctx := v8go.NewContext(nil)
	iso := ctx.Isolate()
	defer iso.Dispose()
	defer ctx.Close()

	cpuProfiler := v8go.NewCPUProfiler(iso)
	defer cpuProfiler.Dispose()

	cpuProfiler.StartProfiling("cpuprofilertest")

	_, err := ctx.RunScript(profileScript, "script.js")
	fatalIf(t, err)
	val, err := ctx.Global().Get("start")
	fatalIf(t, err)
	fn, err := val.AsFunction()
	fatalIf(t, err)
	_, err = fn.Call(ctx.Global())
	fatalIf(t, err)

	cpuProfile := cpuProfiler.StopProfiling("cpuprofilertest")
	if cpuProfile == nil {
		t.Fatal("expected profiler not to be nil")
	}
	root := cpuProfile.GetTopDownRoot()
	if root == nil {
		t.Fatal("expected root not to be nil")
	}
	if root.GetFunctionName() != "(root)" {
		t.Errorf("expected (root), but got %v", root.GetFunctionName())
	}
}

const profileScript = `function loop(timeout) {
  this.mmm = 0;
  var start = Date.now();
  while (Date.now() - start < timeout) {
    var n = 100;
    while(n > 1) {
      n--;
      this.mmm += n * n * n;
    }
  }
}
function delay() { try { loop(10); } catch(e) { } }
function bar() { delay(); }
function baz() { delay(); }
function foo() {
    try {
       delay();
       bar();
       delay();
       baz();
    } catch (e) { }
}
function start(timeout) {
  var start = Date.now();
  do {
    foo();
    var duration = Date.now() - start;
  } while (duration < timeout);
  return duration;
};`
