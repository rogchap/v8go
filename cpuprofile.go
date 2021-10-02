// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import "time"

// CPUProfile contains a CPU profile in a form of top-down call tree
// (from main() down to functions that do all the work).
type CPUProfile struct {
	ptr C.CpuProfilePtr
	iso *Isolate

	title       string
	topDownRoot *CPUProfileNode
	startTime   time.Time
	endTime     time.Time
}

// Returns the root node of the top down call tree.
func (c *CPUProfile) GetTopDownRoot() *CPUProfileNode {
	return c.topDownRoot
}

// Returns CPU profile title.
func (c *CPUProfile) GetTitle() string {
	return c.title
}

// Returns time when the profile recording was started (in microseconds)
// since some unspecified starting point.
func (c *CPUProfile) GetStartTime() time.Time {
	return c.startTime
}

// Returns time when the profile recording was stopped (in microseconds)
// since some unspecified starting point.
// The point is equal to the starting point used by GetStartTime.
func (c *CPUProfile) GetEndTime() time.Time {
	return c.endTime
}

func (c *CPUProfile) Delete() {
	if c.ptr == nil {
		return
	}
	C.CpuProfileDelete(c.ptr)
	c.ptr = nil
}
