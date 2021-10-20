package v8go_test

import (
	"testing"

	v8 "rogchap.com/v8go"
)

func TestCPUProfile(t *testing.T) {
	t.Parallel()

	ctx := v8.NewContext(nil)
	iso := ctx.Isolate()
	defer iso.Dispose()
	defer ctx.Close()

	cpuProfiler := v8.NewCPUProfiler(iso)
	defer cpuProfiler.Dispose()

	title := "cpuprofiletest"
	cpuProfiler.StartProfiling(title)

	_, err := ctx.RunScript(profileScript, "script.js")
	fatalIf(t, err)
	val, err := ctx.Global().Get("start")
	fatalIf(t, err)
	fn, err := val.AsFunction()
	fatalIf(t, err)
	_, err = fn.Call(ctx.Global())
	fatalIf(t, err)

	cpuProfile := cpuProfiler.StopProfiling(title)
	defer cpuProfile.Delete()

	if cpuProfile.GetTitle() != title {
		t.Fatalf("expected title %s, but got %s", title, cpuProfile.GetTitle())
	}

	root := cpuProfile.GetTopDownRoot()
	if root == nil {
		t.Fatal("expected root not to be nil")
	}
	if root.GetFunctionName() != "(root)" {
		t.Errorf("expected (root), but got %v", root.GetFunctionName())
	}

	if !cpuProfile.GetStartTime().Before(cpuProfile.GetEndTime()) {
		t.Fatalf("expected profile start time (%s) before end time (%s)", cpuProfile.GetStartTime(), cpuProfile.GetEndTime())
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
