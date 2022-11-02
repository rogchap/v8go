// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

/*
#include <stdlib.h>
#include "v8go.h"
*/
import "C"
import (
	"time"
	"unsafe"
)

type CPUProfiler struct {
	p   *C.CPUProfiler
	iso *Isolate
}

// CPUProfiler is used to control CPU profiling.
func NewCPUProfiler(iso *Isolate) *CPUProfiler {
	profiler := C.NewCPUProfiler(iso.ptr)
	return &CPUProfiler{
		p:   profiler,
		iso: iso,
	}
}

// Dispose will dispose the profiler.
func (c *CPUProfiler) Dispose() {
	if c.p == nil {
		return
	}

	C.CPUProfilerDispose(c.p)
	c.p = nil
}

// StartProfiling starts collecting a CPU profile. Title may be an empty string. Several
// profiles may be collected at once. Attempts to start collecting several
// profiles with the same title are silently ignored.
func (c *CPUProfiler) StartProfiling(title string) {
	if c.p == nil || c.iso.ptr == nil {
		panic("profiler or isolate are nil")
	}

	tstr := C.CString(title)
	defer C.free(unsafe.Pointer(tstr))

	C.CPUProfilerStartProfiling(c.p, tstr)
}

// Stops collecting CPU profile with a given title and returns it.
// If the title given is empty, finishes the last profile started.
func (c *CPUProfiler) StopProfiling(title string) *CPUProfile {
	if c.p == nil || c.iso.ptr == nil {
		panic("profiler or isolate are nil")
	}

	tstr := C.CString(title)
	defer C.free(unsafe.Pointer(tstr))

	profile := C.CPUProfilerStopProfiling(c.p, tstr)

	return &CPUProfile{
		p:               profile,
		title:           C.GoString(profile.title),
		root:            newCPUProfileNode(profile.root, nil),
		startTimeOffset: time.Duration(profile.startTime) * time.Millisecond,
		endTimeOffset:   time.Duration(profile.endTime) * time.Millisecond,
	}
}

func newCPUProfileNode(node *C.CPUProfileNode, parent *CPUProfileNode) *CPUProfileNode {
	n := &CPUProfileNode{
		nodeId:             int(node.nodeId),
		scriptId:           int(node.scriptId),
		scriptResourceName: C.GoString(node.scriptResourceName),
		functionName:       C.GoString(node.functionName),
		lineNumber:         int(node.lineNumber),
		columnNumber:       int(node.columnNumber),
		hitCount:           int(node.hitCount),
		bailoutReason:      C.GoString(node.bailoutReason),
		parent:             parent,
	}

	if node.childrenCount > 0 {
		n.children = make([]*CPUProfileNode, node.childrenCount)
		for i, child := range (*[1 << 28]*C.CPUProfileNode)(unsafe.Pointer(node.children))[:node.childrenCount:node.childrenCount] {
			n.children[i] = newCPUProfileNode(child, n)
		}
	}

	return n
}
