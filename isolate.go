// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"unsafe"
)

var v8once sync.Once

// ErrIsolateDisposed means isolate has been terminated
var ErrIsolateDisposed = errors.New("v8go: isolate disposed")

// ErrIsolateInUse means isolate is currently in use and cannot be disposed
var ErrIsolateInUse = errors.New("v8go: isolate in use")

// Isolate is a JavaScript VM instance with its own heap and
// garbage collector. Most applications will create one isolate
// with many V8 contexts for execution.
type Isolate struct {
	ptr C.IsolatePtr

	cbMutex *sync.RWMutex
	cbSeq   int
	cbs     map[int]FunctionCallback
	ctx     context.Context
	emux    *sync.Mutex
}

// HeapStatistics represents V8 isolate heap statistics
type HeapStatistics struct {
	TotalHeapSize            uint64
	TotalHeapSizeExecutable  uint64
	TotalPhysicalSize        uint64
	TotalAvailableSize       uint64
	UsedHeapSize             uint64
	HeapSizeLimit            uint64
	MallocedMemory           uint64
	ExternalMemory           uint64
	PeakMallocedMemory       uint64
	NumberOfNativeContexts   uint64
	NumberOfDetachedContexts uint64
}

// NewIsolate creates a new V8 isolate. Only one thread may access
// a given isolate at a time, but different threads may access
// different isolates simultaneously.
// When an isolate is no longer used its resources should be freed
// by calling iso.Dispose().
// An *Isolate can be used as a v8go.ContextOption to create a new
// Context, rather than creating a new default Isolate.
func NewIsolate() (*Isolate, error) {
	return NewIsolateContext(context.Background())
}

// NewIsolateContext creates a new V8 isolate with context. Only one thread may access
// a given isolate at a time, but different threads may access
// different isolates simultaneously.
// When an isolate is no longer used its resources should be freed
// by calling iso.Dispose().
// An *Isolate can be used as a v8go.ContextOption to create a new
// Context, rather than creating a new default Isolate.
func NewIsolateContext(ctx context.Context) (*Isolate, error) {
	v8once.Do(func() {
		C.Init()
	})
	iso := &Isolate{
		ptr:     C.NewIsolate(),
		cbs:     make(map[int]FunctionCallback),
		ctx:     ctx,
		cbMutex: &sync.RWMutex{},
		emux:    &sync.Mutex{},
	}
	// TODO: [RC] catch any C++ exceptions and return as error
	return iso, nil
}

// WithContext adds context to isolate.
func (i *Isolate) WithContext(ctx context.Context) *Isolate {
	iso := new(Isolate)
	*iso = *i
	iso.ctx = ctx
	return iso
}

// Context returns isolate context.
func (i *Isolate) Context() context.Context {
	if i.ctx != nil {
		return i.ctx
	}
	return context.Background()
}

// TerminateExecution terminates forcefully the current thread
// of JavaScript execution in the given isolate.
func (i *Isolate) TerminateExecution() {
	if i.ptr == nil {
		return
	}
	C.IsolateTerminateExecution(i.ptr)
}

// TerminateExecutionWithLock terminates forcefully the current thread
// and aquires lock. This is useful when you need to wait until
// it is fully terminated. There is a possibility that if multiple go routines
// runs into dead-lock.
func (i *Isolate) TerminateExecutionWithLock() {
	if i.ptr == nil {
		return
	}
	C.IsolateTerminateExecution(i.ptr)
	C.IsolateAcquireLock(i.ptr)
}

// IsolateTerminateExecution returns if is V8 terminating JavaScript execution.
// Returns true if JavaScript execution is currently terminating
// because of a call to TerminateExecution. In that case there are
// still JavaScript frames on the stack and the termination
// exception is still active.
func (i *Isolate) IsExecutionTerminating() bool {
	if i.ptr == nil {
		return false
	}
	return C.IsolateIsExecutionTerminating(i.ptr) != 0
}

// IsInUse checks if this isolate is in use.
// True if at least one thread Enter'ed this isolate.
// This will block if any context entered.
func (i *Isolate) IsInUse() bool {
	if i.ptr == nil {
		return false
	}
	return C.IsolateIsInUse(i.ptr) != 0
}

// GetHeapStatistics returns heap statistics for an isolate.
func (i *Isolate) GetHeapStatistics() HeapStatistics {
	hs := C.IsolationGetHeapStatistics(i.ptr)

	return HeapStatistics{
		TotalHeapSize:            uint64(hs.total_heap_size),
		TotalHeapSizeExecutable:  uint64(hs.total_heap_size_executable),
		TotalPhysicalSize:        uint64(hs.total_physical_size),
		TotalAvailableSize:       uint64(hs.total_available_size),
		UsedHeapSize:             uint64(hs.used_heap_size),
		HeapSizeLimit:            uint64(hs.heap_size_limit),
		MallocedMemory:           uint64(hs.malloced_memory),
		ExternalMemory:           uint64(hs.external_memory),
		PeakMallocedMemory:       uint64(hs.peak_malloced_memory),
		NumberOfNativeContexts:   uint64(hs.number_of_native_contexts),
		NumberOfDetachedContexts: uint64(hs.number_of_detached_contexts),
	}
}

// Dispose will dispose the Isolate VM.
func (i *Isolate) Dispose() error {
	if i.ptr == nil {
		return ErrIsolateDisposed
	}
	// In case of Exec method we want to lock
	// dispose. Safe() should be safe to execute and
	// should always finish, if you wish to interrupt safe
	// call TerminateExecution() first
	i.emux.Lock()
	defer i.emux.Unlock()

	// If isolate is in use, aka long running script
	// it should not be allowed to be dispossed as
	// that would cause SIGILL
	if i.IsInUse() {
		return ErrIsolateInUse
	}

	// we do not want to dispose entered context.
	C.IsolateDispose(i.ptr)
	i.ptr = nil
	return nil
}

// Throw an exception into javascript land from within a go function callback
func (i *Isolate) ThrowException(msg string) error {
	if i.ptr == nil {
		return ErrIsolateDisposed
	}
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))
	C.ThrowException(i.ptr, cmsg)
	return nil
}

// Exec will spawn goroutine and run given callback safely.
// While goroutine is running we lock routine to os thread.
// This is necessary for long running scripts as goroutines
// can be moved from one thread to another.
func (i *Isolate) Exec(fn func(*Isolate) error) error {
	ch := make(chan error, 1)
	ctx, cancel := context.WithCancel(i.Context())
	defer cancel()

	i.emux.Lock()
	defer i.emux.Unlock()

	// If isolate is dipised return
	if i.ptr == nil {
		return ErrIsolateDisposed
	}

	// if context has been canceled return
	if ctx.Err() != nil {
		return ctx.Err()
	}
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		ch <- fn(i)
	}()

	return <-ch
}

func (i *Isolate) apply(opts *execContextOptions) {
	opts.iso = i
}

func (i *Isolate) registerCallback(cb FunctionCallback) int {
	i.cbMutex.Lock()
	i.cbSeq++
	ref := i.cbSeq
	i.cbs[ref] = cb
	i.cbMutex.Unlock()
	return ref
}

func (i *Isolate) getCallback(ref int) FunctionCallback {
	i.cbMutex.RLock()
	defer i.cbMutex.RUnlock()
	return i.cbs[ref]
}
