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

	cpuProfile := cpuProfiler.StopProfiling(title)
	if cpuProfile == nil {
		t.Fatal("expected profile not to be nil")
	}
	defer cpuProfile.Delete()

	rootNode := cpuProfile.GetTopDownRoot()
	if rootNode == nil {
		t.Fatal("expected top down root not to be nil")
	}
	checkChildren(t, rootNode, []string{"(program)", "start", "(garbage collector)"})

	startNode := rootNode.GetChild(1)
	checkChildren(t, startNode, []string{"foo"})
	checkNode(t, startNode, "script.js", "start", 23, 15)

	parentName := startNode.GetParent().GetFunctionName()
	if parentName != "(root)" {
		t.Fatalf("expected (root), but got %v", parentName)
	}

	fooNode := startNode.GetChild(0)
	checkChildren(t, fooNode, []string{"delay", "bar", "baz"})
	checkNode(t, fooNode, "script.js", "foo", 15, 13)

	delayNode := fooNode.GetChild(0)
	checkChildren(t, delayNode, []string{"loop"})
	checkNode(t, delayNode, "script.js", "delay", 12, 15)

	barNode := fooNode.GetChild(1)
	checkChildren(t, barNode, []string{"delay"})

	bazNode := fooNode.GetChild(2)
	checkChildren(t, bazNode, []string{"delay"})
}

func checkChildren(t *testing.T, node *v8.CPUProfileNode, names []string) {
	t.Helper()

	for i, n := range names {
		childFunctionName := node.GetChild(i).GetFunctionName()
		if childFunctionName != n {
			t.Fatalf("expected child %d to have name %s, but has %s", i, n, childFunctionName)
		}
	}
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
