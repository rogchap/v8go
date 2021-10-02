// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"fmt"
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
	// t.Parallel()

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
	defer cpuProfile.Delete()

	root := cpuProfile.GetTopDownRoot()
	if root == nil {
		t.Fatal("expected root not to be nil")
	}
	err = checkNode(root, "(root)", 0, 0)
	fatalIf(t, err)
	err = checkChildren(root, []string{"(program)", "start", "(garbage collector)"})
	fatalIf(t, err)

	start := root.GetChild(1)
	err = checkNode(start, "start", 23, 15)
	fatalIf(t, err)
	err = checkChildren(start, []string{"foo"})
	fatalIf(t, err)

	foo := start.GetChild(0)
	err = checkNode(foo, "foo", 15, 13)
	fatalIf(t, err)
	err = checkChildren(foo, []string{"delay", "bar", "baz"})
	fatalIf(t, err)

	baz := foo.GetChild(2)
	err = checkNode(baz, "baz", 14, 13)
	fatalIf(t, err)
	err = checkChildren(baz, []string{"delay"})
	fatalIf(t, err)

	delay := baz.GetChild(0)
	err = checkNode(delay, "delay", 12, 15)
	fatalIf(t, err)
	err = checkChildren(delay, []string{"loop"})
	fatalIf(t, err)
}

func checkChildren(node *v8go.CPUProfileNode, names []string) error {
	nodeName := node.GetFunctionName()
	if node.GetChildrenCount() != len(names) {
		present := []string{}
		for i := 0; i < node.GetChildrenCount(); i++ {
			present = append(present, node.GetChild(i).GetFunctionName())
		}
		return fmt.Errorf("child count for node %s should be %d but was %d: %v", nodeName, len(names), node.GetChildrenCount(), present)
	}
	for i, n := range names {
		if node.GetChild(i).GetFunctionName() != n {
			return fmt.Errorf("expected %s child %d to have name %s", nodeName, i, n)
		}
	}
	return nil
}

func checkNode(node *v8go.CPUProfileNode, name string, line, column int) error {
	if node.GetFunctionName() != name {
		return fmt.Errorf("expected node to have function name `%s` but had `%s`", name, node.GetFunctionName())
	}
	if node.GetLineNumber() != line {
		return fmt.Errorf("expected node %s at line %d, but got %d", name, line, node.GetLineNumber())
	}
	if node.GetColumnNumber() != column {
		return fmt.Errorf("expected node %s at column %d, but got %d", name, column, node.GetColumnNumber())
	}
	return nil
}

// const profileTree = `
// [Top down]:
//  1062     0   (root) [-1]
//  1054     0    start [-1]
//  1054     1      foo [-1]
//   265     0        baz [-1]
//   265     1          delay [-1]
//   264   264            loop [-1]
//   525     3        delay [-1]
//   522   522          loop [-1]
//   263     0        bar [-1]
//   263     1          delay [-1]
//   262   262            loop [-1]
//     2     2    (program) [-1]
//     6     6    (garbage collector) [-1]
// `

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
