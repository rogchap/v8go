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

	// startTimeOffset is the time when the profile recording was started
	// since some unspecified starting point.
	startTimeOffset time.Duration

	// endTimeOffset is the time when the profile recording was stopped
	// since some unspecified starting point.
	// The point is equal to the starting point used by startTimeOffset.
	endTimeOffset time.Duration
}

// Returns CPU profile title.
func (c *CPUProfile) GetTitle() string {
	return c.title
}

// Returns the root node of the top down call tree.
func (c *CPUProfile) GetTopDownRoot() *CPUProfileNode {
	return c.root
}

// Returns the duration of the profile.
func (c *CPUProfile) GetDuration() time.Duration {
	return c.endTimeOffset - c.startTimeOffset
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
