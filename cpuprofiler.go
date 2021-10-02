// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"time"
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

// Dispose will dispose the profiler.
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
// TODO: Return CPUProfilingStats
// Enum for returning profiling status. Once StartProfiling is called,
// we want to return to clients whether the profiling was able to start
// correctly, or return a descriptive error.
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

	profilePtr := C.CpuProfilerStopProfiling(c.iso.ptr, c.ptr, tstr)

	rootPtr := C.CpuProfileGetTopDownRoot(profilePtr)
	rootNode := &CPUProfileNode{
		parent:             nil,
		scriptResourceName: C.GoString(C.CpuProfileNodeGetScriptResourceName(rootPtr)),
		functionName:       C.GoString(C.CpuProfileNodeGetFunctionName(rootPtr)),
		lineNumber:         int(C.CpuProfileNodeGetLineNumber(rootPtr)),
		columnNumber:       int(C.CpuProfileNodeGetColumnNumber(rootPtr)),
	}
	rootNode.children = getChildren(rootNode, rootPtr)

	return &CPUProfile{
		ptr:         profilePtr,
		iso:         c.iso,
		title:       title,
		topDownRoot: rootNode,
		startTime:   time.Unix(0, int64(C.CpuProfileGetStartTime(profilePtr))/1000),
		endTime:     time.Unix(0, int64(C.CpuProfileGetEndTime(profilePtr))/1000),
	}
}

func getChildren(parent *CPUProfileNode, ptr C.CpuProfileNodePtr) []*CPUProfileNode {
	count := C.CpuProfileNodeGetChildrenCount(ptr)
	if int(count) == 0 {
		return []*CPUProfileNode{}
	}

	parent.children = make([]*CPUProfileNode, count)

	for i := 0; i < int(count); i++ {
		childNodePtr := C.CpuProfileNodeGetChild(ptr, C.int(i))

		childNode := &CPUProfileNode{
			parent:             parent,
			scriptResourceName: C.GoString(C.CpuProfileNodeGetScriptResourceName(childNodePtr)),
			functionName:       C.GoString(C.CpuProfileNodeGetFunctionName(childNodePtr)),
			lineNumber:         int(C.CpuProfileNodeGetLineNumber(childNodePtr)),
			columnNumber:       int(C.CpuProfileNodeGetColumnNumber(childNodePtr)),
		}
		childNode.children = getChildren(childNode, childNodePtr)
		parent.children[i] = childNode
	}
	return parent.children
}
