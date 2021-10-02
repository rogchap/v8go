// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
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

	cpuProfiler.StartProfiling("cpuprofilenodetest")

	_, _ = ctx.RunScript(profileScript, "script.js")
	val, _ := ctx.Global().Get("start")
	fn, _ := val.AsFunction()
	_, _ = fn.Call(ctx.Global())

	cpuProfile := cpuProfiler.StopProfiling("cpuprofilenodetest")
	defer cpuProfile.Delete()

	node := cpuProfile.GetTopDownRoot()
	err := checkNode(node, "(root)", 0, 0)
	fatalIf(t, err)

	if node.GetParent() != nil {
		t.Fatal("expected root node to have nil parent")
	}

	if node.GetChildrenCount() < 2 {
		t.Fatalf("expected at least 2 children, but got %d", node.GetChildrenCount())
	}

	if node.GetChild(1).GetFunctionName() != "start" {
		t.Fatalf("expected child node with name `start` but got %s", node.GetChild(1).GetFunctionName())
	}

	if node.GetChild(1).GetScriptResourceName() != "script.js" {
		t.Fatalf("expected child to have script resource name `script.js` but had `%s`", node.GetScriptResourceName())
	}

	if node.GetChild(0).GetParent() != node {
		t.Fatal("expected child's parent to be the same node")
	}
}
