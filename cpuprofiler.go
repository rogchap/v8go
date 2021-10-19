// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
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
		p:         profile,
		Title:     C.GoString(profile.title),
		Root:      NewCPUProfileNode(profile.root, nil),
		StartTime: time.Unix(0, int64(profile.startTime)*1000),
		EndTime:   time.Unix(0, int64(profile.endTime)*1000),
	}
}

type CPUProfile struct {
	p *C.CPUProfile

	Title     string
	Root      *CPUProfileNode
	StartTime time.Time
	EndTime   time.Time
}

type CPUProfileNode struct {
	p *C.CPUProfileNode

	ScriptResourceName string
	FunctionName       string
	LineNumber         int
	ColumnNumber       int
	Children           []*CPUProfileNode
	Parent             *CPUProfileNode
}

func NewCPUProfileNode(node *C.CPUProfileNode, parent *CPUProfileNode) *CPUProfileNode {
	n := &CPUProfileNode{
		p:                  node,
		ScriptResourceName: C.GoString(node.scriptResourceName),
		FunctionName:       C.GoString(node.functionName),
		LineNumber:         int(node.lineNumber),
		ColumnNumber:       int(node.columnNumber),
		Parent:             parent,
	}

	if node.childrenCount > 0 {
		for _, child := range (*[1 << 28]*C.CPUProfileNode)(unsafe.Pointer(node.children))[:node.childrenCount:node.childrenCount] {
			n.Children = append(n.Children, NewCPUProfileNode(child, n))
		}
	}

	return n
}

func (c *CPUProfile) Delete() {
	if c.p == nil {
		return
	}
	C.CPUProfileDelete(c.p)
	c.p = nil
}
