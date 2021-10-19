// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	v8 "rogchap.com/v8go"
)

func TestCPUProfiler_Dispose(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	defer iso.Dispose()
	cpuProfiler := v8.NewCPUProfiler(iso)

	cpuProfiler.Dispose()
	// noop when called multiple times
	cpuProfiler.Dispose()

	// verify panics when profiler disposed
	if recoverPanic(func() { cpuProfiler.StartProfiling("") }) == nil {
		t.Error("expected panic")
	}

	if recoverPanic(func() { cpuProfiler.StopProfiling("") }) == nil {
		t.Error("expected panic")
	}

	cpuProfiler = v8.NewCPUProfiler(iso)
	defer cpuProfiler.Dispose()
	iso.Dispose()

	// verify panics when isolate disposed
	if recoverPanic(func() { cpuProfiler.StartProfiling("") }) == nil {
		t.Error("expected panic")
	}

	if recoverPanic(func() { cpuProfiler.StopProfiling("") }) == nil {
		t.Error("expected panic")
	}
}

func TestCPUProfiler(t *testing.T) {
	t.Parallel()

	ctx := v8.NewContext(nil)
	iso := ctx.Isolate()
	defer iso.Dispose()
	defer ctx.Close()

	cpuProfiler := v8.NewCPUProfiler(iso)
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
	defer cpuProfile.Delete()

	root := cpuProfile.Root
	if root == nil {
		t.Fatal("expected root not to be nil")
	}
	if root.FunctionName != "(root)" {
		t.Errorf("expected (root), but got %v", root.FunctionName)
	}
}

func TestCPUProfile_Delete(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	defer iso.Dispose()

	cpuProfiler := v8.NewCPUProfiler(iso)
	defer cpuProfiler.Dispose()

	cpuProfiler.StartProfiling("cpuprofiletest")
	cpuProfile := cpuProfiler.StopProfiling("cpuprofiletest")
	cpuProfile.Delete()
	// noop when called multiple times
	cpuProfile.Delete()
}

func TestCPUProfile(t *testing.T) {
	t.Parallel()

	ctx := v8.NewContext(nil)
	iso := ctx.Isolate()
	defer iso.Dispose()
	defer ctx.Close()

	cpuProfiler := v8.NewCPUProfiler(iso)
	defer cpuProfiler.Dispose()

	cpuProfiler.StartProfiling("cpuprofiletest")

	_, err := ctx.RunScript(profileScript, "script.js")
	fatalIf(t, err)
	val, err := ctx.Global().Get("start")
	fatalIf(t, err)
	fn, err := val.AsFunction()
	fatalIf(t, err)
	_, err = fn.Call(ctx.Global())
	fatalIf(t, err)

	cpuProfile := cpuProfiler.StopProfiling("cpuprofiletest")
	defer cpuProfile.Delete()

	if cpuProfile.Title != "cpuprofiletest" {
		t.Errorf("expected cpuprofiletest, but got %v", cpuProfile.Title)
	}

	root := cpuProfile.Root
	if root == nil {
		t.Fatal("expected root not to be nil")
	}
	if root.FunctionName != "(root)" {
		t.Errorf("expected (root), but got %v", root.FunctionName)
	}
}

func TestCPUProfileNode(t *testing.T) {
	t.Parallel()

	ctx := v8.NewContext(nil)
	iso := ctx.Isolate()
	defer iso.Dispose()
	defer ctx.Close()

	cpuProfiler := v8.NewCPUProfiler(iso)
	defer cpuProfiler.Dispose()

	cpuProfiler.StartProfiling("cpuprofilenodetest")

	_, err := ctx.RunScript(profileScript, "script.js")
	fatalIf(t, err)
	val, err := ctx.Global().Get("start")
	fatalIf(t, err)
	fn, err := val.AsFunction()
	fatalIf(t, err)
	timeout, err := v8.NewValue(iso, int32(1000))
	fatalIf(t, err)
	_, err = fn.Call(ctx.Global(), timeout)
	fatalIf(t, err)

	cpuProfile := cpuProfiler.StopProfiling("cpuprofilenodetest")
	if cpuProfile == nil {
		t.Fatal("expected profile not to be nil")
	}
	defer cpuProfile.Delete()

	if !cpuProfile.StartTime.Before(cpuProfile.EndTime) {
		t.Fatal("expected start time before end time")
	}

	rootNode := cpuProfile.Root
	if rootNode == nil {
		t.Fatal("expected top down root not to be nil")
	}
	if rootNode.FunctionName != "(root)" {
		t.Fatalf("expected (root), but got %v", rootNode.FunctionName)
	}

	checkChildren(t, rootNode, []string{"(program)", "start", "(garbage collector)"})

	if len(rootNode.Children) != 3 {
		t.Fatalf("expected 3 children, but got %d", len(rootNode.Children))
	}

	startNode := rootNode.Children[1]
	checkChildren(t, startNode, []string{"foo"})
	checkNode(t, startNode, "script.js", "start", 23, 15)

	parentName := startNode.Parent.FunctionName
	if parentName != "(root)" {
		t.Fatalf("expected (root), but got %v", parentName)
	}

	fooNode := startNode.Children[0]
	checkChildren(t, fooNode, []string{"delay", "bar", "baz"})
	checkNode(t, fooNode, "script.js", "foo", 15, 13)

	delayNode := fooNode.Children[0]
	checkChildren(t, delayNode, []string{"loop"})
	checkNode(t, delayNode, "script.js", "delay", 12, 15)

	barNode := fooNode.Children[1]
	checkChildren(t, barNode, []string{"delay"})

	bazNode := fooNode.Children[2]
	checkChildren(t, bazNode, []string{"delay"})
}

func checkChildren(t *testing.T, node *v8.CPUProfileNode, names []string) {
	t.Helper()

	nodeName := node.FunctionName
	if len(node.Children) != len(names) {
		t.Fatalf("expected %d children for node %s, found %d: %#v", len(names), nodeName, len(node.Children), node.Children)
	}

	for i, n := range names {
		if node.Children[i].FunctionName != n {
			t.Fatalf("expected %s child %d to have name %s, but has %s: %#v", nodeName, i, n, node.Children[i].FunctionName, node.Children)
		}
	}
}

func checkNode(t *testing.T, node *v8.CPUProfileNode, scriptResourceName string, functionName string, line, column int) {
	t.Helper()

	if node.FunctionName != functionName {
		t.Fatalf("expected node to have function name %s, but got %s", functionName, node.FunctionName)
	}
	if node.ScriptResourceName != scriptResourceName {
		t.Fatalf("expected node to have script resource name %s, but got %s", scriptResourceName, node.ScriptResourceName)
	}
	if node.LineNumber != line {
		t.Fatalf("expected node at line %d, but got %d", line, node.LineNumber)
	}
	if node.ColumnNumber != column {
		t.Fatalf("expected node at column %d, but got %d", column, node.ColumnNumber)
	}
}

const profileScript = `function loop(timeout) {
  this.mmm = 0;
  var start = Date.now();
  while (Date.now() - start < timeout) {
    var n = 10;
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
