// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

type CPUProfileNode struct {
	// The id of the current node, unique within the tree.
	nodeId int

	// The id of the script where the function originates.
	scriptId int

	// The resource name for script from where the function originates.
	scriptResourceName string

	// The function name (empty string for anonymous functions.)
	functionName string

	// The number of the line where the function originates.
	lineNumber int

	// The number of the column where the function originates.
	columnNumber int

	// The count of samples where the function was currently executing.
	hitCount int

	// The bailout reason for the function if the optimization was disabled for it.
	bailoutReason string

	// The children node of this node.
	children []*CPUProfileNode

	// The parent node of this node.
	parent *CPUProfileNode
}

// Returns node id.
func (c *CPUProfileNode) GetNodeId() int {
	return c.nodeId
}

// Returns id for script from where the function originates.
func (c *CPUProfileNode) GetScriptId() int {
	return c.scriptId
}

// Returns function name (empty string for anonymous functions.)
func (c *CPUProfileNode) GetFunctionName() string {
	return c.functionName
}

// Returns resource name for script from where the function originates.
func (c *CPUProfileNode) GetScriptResourceName() string {
	return c.scriptResourceName
}

// Returns number of the line where the function originates.
func (c *CPUProfileNode) GetLineNumber() int {
	return c.lineNumber
}

// Returns number of the column where the function originates.
func (c *CPUProfileNode) GetColumnNumber() int {
	return c.columnNumber
}

// Returns count of samples where the function was currently executing.
func (c *CPUProfileNode) GetHitCount() int {
	return c.hitCount
}

// Returns the bailout reason for the function if the optimization was disabled for it.
func (c *CPUProfileNode) GetBailoutReason() string {
	return c.bailoutReason
}

// Retrieves the ancestor node, or nil if the root.
func (c *CPUProfileNode) GetParent() *CPUProfileNode {
	return c.parent
}

func (c *CPUProfileNode) GetChildrenCount() int {
	return len(c.children)
}

// Retrieves a child node by index.
func (c *CPUProfileNode) GetChild(index int) *CPUProfileNode {
	return c.children[index]
}
