// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
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

func (c *CPUProfiler) StartProfiling(title string) {
	tstr := C.CString(title)
	defer C.free(unsafe.Pointer(tstr))
	C.CpuProfilerStartProfiling(c.iso.ptr, c.ptr, tstr)
}

func (c *CPUProfiler) StopProfiling(title string, securityToken string) *CPUProfile {
	tstr := C.CString(title)
	defer C.free(unsafe.Pointer(tstr))
	ptr := C.CpuProfilerStopProfiling(c.iso.ptr, c.ptr, tstr)
	return &CPUProfile{ptr: ptr, iso: c.iso}
}
