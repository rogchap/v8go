// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"unsafe"
)

type CPUProfiler struct {
	ptr C.CpuProfilerPtr
	iso *Isolate
}

// CPUProfiler is used to control CPU profiling.
func NewCPUProfiler(iso *Isolate) *CPUProfiler {
	return &CPUProfiler{
		ptr: C.NewCpuProfiler(iso.ptr),
		iso: iso,
	}
}

// Dispose will dispose the profiler; subsequent calls will panic.
func (c *CPUProfiler) Dispose() {
	if c.ptr == nil {
		return
	}
	C.CpuProfilerDispose(c.ptr)
	c.ptr = nil
}

// StartProfiling starts collecting a CPU profile. Title may be an empty string. Several
// profiles may be collected at once. Attempts to start collecting several
// profiles with the same title are silently ignored.
func (c *CPUProfiler) StartProfiling(title string) {
	if c.ptr == nil || c.iso.ptr == nil {
		return
	}
	tstr := C.CString(title)
	defer C.free(unsafe.Pointer(tstr))
	C.CpuProfilerStartProfiling(c.iso.ptr, c.ptr, tstr)
}

// Stops collecting CPU profile with a given title and returns it.
// If the title given is empty, finishes the last profile started.
func (c *CPUProfiler) StopProfiling(title string) *CPUProfile {
	if c.ptr == nil || c.iso.ptr == nil {
		return nil
	}
	tstr := C.CString(title)
	defer C.free(unsafe.Pointer(tstr))
	return &CPUProfile{
		ptr: C.CpuProfilerStopProfiling(c.iso.ptr, c.ptr, tstr),
		iso: c.iso,
	}
}
