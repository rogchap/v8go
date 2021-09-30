// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"

// CPUProfileNode represents a node in a call graph.
type CPUProfileNode struct {
	ptr C.CpuProfileNodePtr
	iso *Isolate
}

// Returns function name (empty string for anonymous functions.)
func (c *CPUProfileNode) GetFunctionName() string {
	str := C.CpuProfileNodeGetFunctionName(c.ptr)
	return C.GoString(str)
}

// Retrieves number of children.
func (c *CPUProfileNode) GetChildrenCount() int {
	i := C.CpuProfileNodeGetChildrenCount(c.ptr)
	return int(i)
}

// Retrieves a child node by index.
func (c *CPUProfileNode) GetChild(index int) *CPUProfileNode {
	count := c.GetChildrenCount()
	if index < 0 || index > count {
		return nil
	}
	ptr := C.CpuProfileNodeGetChild(c.ptr, C.int(index))
	return &CPUProfileNode{ptr: ptr, iso: c.iso}
}

// Retrieves the ancestor node, or null if the root.
func (c *CPUProfileNode) GetParent() *CPUProfileNode {
	ptr := C.CpuProfileNodeGetParent(c.ptr)
	return &CPUProfileNode{ptr: ptr, iso: c.iso}
}

// Returns the number, 1-based, of the line where the function originates.
func (c *CPUProfileNode) GetLineNumber() int {
	i := C.CpuProfileNodeGetLineNumber(c.ptr)
	return int(i)
}

//  Returns 1-based number of the column where the function originates.
func (c *CPUProfileNode) GetColumnNumber() int {
	i := C.CpuProfileNodeGetColumnNumber(c.ptr)
	return int(i)
}
