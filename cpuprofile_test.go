// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"
	"time"

	v8 "rogchap.com/v8go"
)

func TestCPUProfile_Dispose(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	defer iso.Dispose()

	cpuProfiler := v8.NewCPUProfiler(iso)
	defer cpuProfiler.Dispose()

	cpuProfiler.StartProfiling("cpuprofiledispose")
	cpuProfile := cpuProfiler.StopProfiling("cpuprofiledispose")
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

	_, err := ctx.RunScript(`function foo() {  }; foo();`, "script.js")
	fatalIf(t, err)

	// Ensure different start/end time
	time.Sleep(10 * time.Microsecond)

	cpuProfile := cpuProfiler.StopProfiling("cpuprofiletest")
	if cpuProfile == nil {
		t.Fatal("expected profiler not to be nil")
	}
	defer cpuProfile.Delete()

	if cpuProfile.GetTitle() != "cpuprofiletest" {
		t.Errorf("expected cpuprofiletest, but got %v", cpuProfile.GetTitle())
	}

	if cpuProfile.GetTopDownRoot() == nil {
		t.Fatal("expected root not to be nil")
	}

	if cpuProfile.GetStartTime().Equal(cpuProfile.GetEndTime()) {
		t.Fatal("expected different start and end times")
	}
}
