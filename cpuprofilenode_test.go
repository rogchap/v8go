// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	v8 "rogchap.com/v8go"
)

func TestCPUProfileNode(t *testing.T) {
	t.Parallel()

	ctx := v8.NewContext(nil)
	iso := ctx.Isolate()
	defer iso.Dispose()
	defer ctx.Close()

	cpuProfiler := v8.NewCPUProfiler(iso)
	defer cpuProfiler.Dispose()

	title := "cpuprofilenodetest"
	cpuProfiler.StartProfiling(title)

	_, err := ctx.CompileAndRun(profileScript, "script.js")
	fatalIf(t, err)
	val, err := ctx.Global().Get("start")
	fatalIf(t, err)
	fn, err := val.AsFunction()
	fatalIf(t, err)
	timeout, err := v8.NewValue(iso, int32(1000))
	fatalIf(t, err)
	_, err = fn.Call(ctx.Global(), timeout)
	fatalIf(t, err)

	cpuProfile := cpuProfiler.StopProfiling(title)
	if cpuProfile == nil {
		t.Fatal("expected profile not to be nil")
	}
	defer cpuProfile.Delete()

	rootNode := cpuProfile.GetTopDownRoot()
	if rootNode == nil {
		t.Fatal("expected top down root not to be nil")
	}
	count := rootNode.GetChildrenCount()
	var startNode *v8.CPUProfileNode
	for i := 0; i < count; i++ {
		if rootNode.GetChild(i).GetFunctionName() == "start" {
			startNode = rootNode.GetChild(i)
		}
	}
	if startNode == nil {
		t.Fatal("expected node not to be nil")
	}
	checkNode(t, startNode, "script.js", "start", 23, 15)

	parentName := startNode.GetParent().GetFunctionName()
	if parentName != "(root)" {
		t.Fatalf("expected (root), but got %v", parentName)
	}

	fooNode := findChild(t, startNode, "foo")
	checkNode(t, fooNode, "script.js", "foo", 15, 13)

	delayNode := findChild(t, fooNode, "delay")
	checkNode(t, delayNode, "script.js", "delay", 12, 15)

	barNode := findChild(t, fooNode, "bar")
	checkNode(t, barNode, "script.js", "bar", 13, 13)

	loopNode := findChild(t, delayNode, "loop")
	checkNode(t, loopNode, "script.js", "loop", 1, 14)

	bazNode := findChild(t, fooNode, "baz")
	checkNode(t, bazNode, "script.js", "baz", 14, 13)
}

func findChild(t *testing.T, node *v8.CPUProfileNode, functionName string) *v8.CPUProfileNode {
	t.Helper()

	var child *v8.CPUProfileNode
	count := node.GetChildrenCount()
	for i := 0; i < count; i++ {
		if node.GetChild(i).GetFunctionName() == functionName {
			child = node.GetChild(i)
		}
	}
	if child == nil {
		t.Fatal("failed to find child node")
	}
	return child
}

func checkNode(t *testing.T, node *v8.CPUProfileNode, scriptResourceName string, functionName string, line, column int) {
	t.Helper()

	if node.GetFunctionName() != functionName {
		t.Fatalf("expected node to have function name %s, but got %s", functionName, node.GetFunctionName())
	}
	if node.GetScriptResourceName() != scriptResourceName {
		t.Fatalf("expected node to have script resource name %s, but got %s", scriptResourceName, node.GetScriptResourceName())
	}
	if node.GetLineNumber() != line {
		t.Fatalf("expected node at line %d, but got %d", line, node.GetLineNumber())
	}
	if node.GetColumnNumber() != column {
		t.Fatalf("expected node at column %d, but got %d", column, node.GetColumnNumber())
	}
}
