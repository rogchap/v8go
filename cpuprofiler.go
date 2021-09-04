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

type CpuProfiler struct {
	ptr C.CpuProfilerPtr
	iso *Isolate
}

// CpuProfile contains a CPU profile in a form of top-down call tree
// (from main() down to functions that do all the work).
type CpuProfile struct {
	ptr C.CpuProfilePtr
	iso *Isolate
}

// CpuProfileNode represents a node in a call graph.
type CpuProfileNode struct {
	ptr C.CpuProfileNodePtr
	iso *Isolate
}

// Returns the root node of the top down call tree.
func (c *CpuProfile) GetTopDownRoot() *CpuProfileNode {
	ptr := C.CpuProfileGetTopDownRoot(c.ptr)
	return &CpuProfileNode{ptr: ptr, iso: c.iso}
}

// Returns function name (empty string for anonymous functions.)
func (c *CpuProfileNode) GetFunctionName() string {
	str := C.CpuProfileNodeGetFunctionName(c.ptr)
	return C.GoString(str)
}

// Retrieves number of children.
func (c *CpuProfileNode) GetChildrenCount() int {
	i := C.CpuProfileNodeGetChildrenCount(c.ptr)
	return int(i)
}

// Retrieves a child node by index.
func (c *CpuProfileNode) GetChild(index int) *CpuProfileNode {
	count := c.GetChildrenCount()
	if index < 0 || index > count {
		return nil
	}
	ptr := C.CpuProfileNodeGetChild(c.ptr, C.int(index))
	return &CpuProfileNode{ptr: ptr, iso: c.iso}
}

// Retrieves the ancestor node, or null if the root.
func (c *CpuProfileNode) GetParent() *CpuProfileNode {
	ptr := C.CpuProfileNodeGetParent(c.ptr)
	return &CpuProfileNode{ptr: ptr, iso: c.iso}
}

// Returns the number, 1-based, of the line where the function originates.
// kNoLineNumberInfo if no line number information is available.
func (c *CpuProfileNode) GetLineNumber() int {
	no := C.CpuProfileNodeGetLineNumber(c.ptr)
	return int(no)
}

//  Returns 1-based number of the column where the function originates.
//  kNoColumnNumberInfo if no column number information is available.
func (c *CpuProfileNode) GetColumnNumber() int {
	no := C.CpuProfileNodeGetColumnNumber(c.ptr)
	return int(no)
}

func NewCpuProfiler(iso *Isolate) *CpuProfiler {
	return &CpuProfiler{
		ptr: C.NewCpuProfiler(iso.ptr),
		iso: iso,
	}
}

func (c *CpuProfiler) StartProfiling(title string) {
	tstr := C.CString(title)
	defer C.free(unsafe.Pointer(tstr))
	C.CpuProfilerStartProfiling(c.iso.ptr, c.ptr, tstr)
}

func (c *CpuProfiler) StopProfiling(title string, securityToken string) *CpuProfile {
	tstr := C.CString(title)
	defer C.free(unsafe.Pointer(tstr))
	ptr := C.CpuProfilerStopProfiling(c.iso.ptr, c.ptr, tstr)
	return &CpuProfile{ptr: ptr, iso: c.iso}
}

// Dispose will dispose the profiler; subsequent calls will panic.
func (c *CpuProfiler) Dispose() {
	if c.ptr == nil {
		return
	}
	C.CpuProfilerDispose(c.ptr)
	c.ptr = nil
}

// TODO
// The profiler must be dispoed after use by calling Dispose()
// static int 	GetProfilesCount ()
// static const CpuProfile * 	GetProfile (int index, Handle< Value > security_token=Handle< Value >())
// static const CpuProfile * 	FindProfile (unsigned uid, Handle< Value > security_token=Handle< Value >())
// static void 	DeleteAllProfiles ()
