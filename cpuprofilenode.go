// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"

// CPUProfileNode represents a node in a call graph.
type CPUProfileNode struct {
	scriptResourceName string
	functionName       string
	lineNumber         int
	columnNumber       int

	parent   *CPUProfileNode
	children []*CPUProfileNode
}

// Retrieves number of children.
func (c *CPUProfileNode) GetChildrenCount() int {
	return len(c.children)
}

// Retrieves a child node by index.
func (c *CPUProfileNode) GetChild(index int) *CPUProfileNode {
	if index < 0 || index > len(c.children) {
		return nil
	}
	return c.children[index]
}

func (c *CPUProfileNode) GetParent() *CPUProfileNode {
	return c.parent
}

func (c *CPUProfileNode) GetScriptResourceName() string {
	return c.scriptResourceName
}

// Returns function name (empty string for anonymous functions.)
func (c *CPUProfileNode) GetFunctionName() string {
	return c.functionName
}

// Returns the number, 1-based, of the line where the function originates.
func (c *CPUProfileNode) GetLineNumber() int {
	return c.lineNumber
}

//  Returns 1-based number of the column where the function originates.
func (c *CPUProfileNode) GetColumnNumber() int {
	return c.columnNumber
}
