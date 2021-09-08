// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	"rogchap.com/v8go"
)

func TestCPUProfileNode(t *testing.T) {
	t.Parallel()

	ctx, _ := v8go.NewContext(nil)
	iso := ctx.Isolate()
	defer iso.Dispose()
	defer ctx.Close()

	cpuProfiler := v8go.NewCPUProfiler(iso)
	defer cpuProfiler.Dispose()

	cpuProfiler.StartProfiling("cpuprofilenodetest")

	_, err := ctx.RunScript(profileScript, "script.js")
	failIf(t, err)
	val, err := ctx.Global().Get("start")
	failIf(t, err)
	fn, err := val.AsFunction()
	failIf(t, err)
	_, err = fn.Call()
	failIf(t, err)

	cpuProfile := cpuProfiler.StopProfiling("cpuprofilenodetest")
	if cpuProfile == nil {
		t.Fatal("expected profile not to be nil")
	}

	root := cpuProfile.GetTopDownRoot()
	if root == nil {
		t.Fatal("expected top down root not to be nil")
	}
	if root.GetFunctionName() != "(root)" {
		t.Errorf("expected (root), but got %v", root.GetFunctionName())
	}
	checkChildren(t, root, []string{"(program)", "start", "(garbage collector)"})

	invalidChild := root.GetChild(4)
	if invalidChild != nil {
		t.Errorf("expected nil child, but got %v", invalidChild.GetFunctionName())
	}

	startNode := root.GetChild(1)
	if startNode.GetFunctionName() != "start" {
		t.Errorf("expected start, but got %v", startNode.GetFunctionName())
	}
	checkChildren(t, startNode, []string{"foo"})
	checkPosition(t, startNode, 23, 15)

	parentName := startNode.GetParent().GetFunctionName()
	if parentName != "(root)" {
		t.Errorf("expected (root), but got %v", parentName)
	}

	fooNode := startNode.GetChild(0)
	checkChildren(t, fooNode, []string{"delay", "bar", "baz"})
	checkPosition(t, fooNode, 15, 13)

	delayNode := fooNode.GetChild(0)
	checkChildren(t, delayNode, []string{"loop"})
	checkPosition(t, delayNode, 12, 15)

	barNode := fooNode.GetChild(1)
	checkChildren(t, barNode, []string{"delay"})

	bazNode := fooNode.GetChild(2)
	checkChildren(t, bazNode, []string{"delay"})
}

func checkChildren(t *testing.T, node *v8go.CPUProfileNode, names []string) {
	nodeName := node.GetFunctionName()
	if node.GetChildrenCount() != len(names) {
		t.Fatalf("expected child count for node %s to equal length of child names", nodeName)
	}
	for i, n := range names {
		if node.GetChild(i).GetFunctionName() != n {
			t.Errorf("expected %s child %d to have name %s", nodeName, i, n)
		}
	}
}

func checkPosition(t *testing.T, node *v8go.CPUProfileNode, line, column int) {
	nodeName := node.GetFunctionName()
	if node.GetLineNumber() != line {
		t.Errorf("expected node %s at line %d, but got %d", nodeName, line, node.GetLineNumber())
	}
	if node.GetColumnNumber() != column {
		t.Errorf("expected node %s at column %d, but got %d", nodeName, column, node.GetColumnNumber())
	}
}
