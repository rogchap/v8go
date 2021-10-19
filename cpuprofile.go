package v8go

/*
#include "v8go.h"
*/
import "C"
import "time"

type CPUProfile struct {
	p *C.CPUProfile

	// The CPU profile title.
	title string

	// root is the root node of the top down call tree.
	root *CPUProfileNode

	// startTime is the time when the profile recording was started (in microseconds)
	// since some unspecified starting point.
	startTime time.Time

	// endTime is the time when the profile recording was stopped (in microseconds)
	// since some unspecified starting point.
	// The point is equal to the starting point used by startTime.
	endTime time.Time
}

// Returns CPU profile title.
func (c *CPUProfile) GetTitle() string {
	return c.title
}

// Returns the root node of the top down call tree.
func (c *CPUProfile) GetTopDownRoot() *CPUProfileNode {
	return c.root
}

// Returns the time when the profile recording was started (in microseconds)
// since some unspecified starting point.
func (c *CPUProfile) GetStartTime() time.Time {
	return c.startTime
}

// Returns the time when the profile recording was stopped (in microseconds)
// since some unspecified starting point.
// The point is equal to the starting point used by startTime.
func (c *CPUProfile) GetEndTime() time.Time {
	return c.endTime
}

// Deletes the profile and removes it from CpuProfiler's list.
// All pointers to nodes previously returned become invalid.
func (c *CPUProfile) Delete() {
	if c.p == nil {
		return
	}
	C.CPUProfileDelete(c.p)
	c.p = nil
}
