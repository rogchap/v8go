package v8go

/*
#include <stdlib.h>
#include "v8go.h"
*/
import "C"
import "unsafe"

type HeapProfiler struct {
	p   *C.V8HeapProfiler
	iso *Isolate
}

func NewHeapProfiler(iso *Isolate) *HeapProfiler {
	profiler := C.NewHeapProfiler(iso.ptr)
	return &HeapProfiler{
		p:   profiler,
		iso: iso,
	}
}

func (c *HeapProfiler) TakeHeapSnapshot() (string, error) {
	if c.p == nil || c.iso.ptr == nil {
		panic("heap profiler or isolate is nil")
	}

	str := C.TakeHeapSnapshot(c.p)
	defer C.free(unsafe.Pointer(str))
	return C.GoString(str), nil
}
