// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"errors"
	"unsafe"
)

type FunctionCodeHandling int

//  Clear - does not keeps any compiled data prior to serialization/deserialization/verify pass
//  Keep - keeps any compiled data prior to serialization/deserialization/verify pass
const (
	FunctionCodeHandlingClear FunctionCodeHandling = iota
	FunctionCodeHandlingKeep
)

// StartupData stores the snapshot blob data
type StartupData struct {
	data     []byte
	raw_size C.int
}

// SnapshotCreator allows creating snapshot.
type SnapshotCreator struct {
	ptr                 C.SnapshotCreatorPtr
	iso                 *Isolate
	defaultContextAdded bool
}

// NewSnapshotCreator creates a new snapshot creator.
func NewSnapshotCreator() *SnapshotCreator {
	v8once.Do(func() {
		C.Init()
	})

	rtn := C.NewSnapshotCreator()

	return &SnapshotCreator{
		ptr:                 rtn.creator,
		iso:                 &Isolate{ptr: rtn.iso},
		defaultContextAdded: false,
	}
}

// GetIsolate returns the Isolate associated with the SnapshotCreator.
// This Isolate must be use to create the contexts that later will be use to create the snapshot blob.
func (s *SnapshotCreator) GetIsolate() (*Isolate, error) {
	if s.ptr == nil {
		return nil, errors.New("v8go: Cannot get Isolate after creating the blob")
	}

	return s.iso, nil
}

// SetDefaultContext set the default context to be included in the snapshot blob.
func (s *SnapshotCreator) SetDefaultContext(ctx *Context) error {
	if s.defaultContextAdded {
		return errors.New("v8go: Cannot set multiple default context for snapshot creator")
	}

	C.SetDefaultContext(s.ptr, ctx.ptr)
	s.defaultContextAdded = true
	ctx.ptr = nil

	return nil
}

// AddContext add additional context to be included in the snapshot blob.
// Returns the index of the context in the snapshot blob, that later can be use to call v8go.NewContextFromSnapshot.
func (s *SnapshotCreator) AddContext(ctx *Context) (int, error) {
	if s.ptr == nil {
		return 0, errors.New("v8go: Cannot add context to snapshot creator after creating the blob")
	}

	index := C.AddContext(s.ptr, ctx.ptr)
	ctx.ptr = nil

	return int(index), nil
}

// Create creates a snapshot data blob.
func (s *SnapshotCreator) Create(functionCode FunctionCodeHandling) (*StartupData, error) {
	if s.ptr == nil {
		return nil, errors.New("v8go: Cannot use snapshot creator after creating the blob")
	}

	if !s.defaultContextAdded {
		return nil, errors.New("v8go: Cannot create a snapshot without a default context")
	}

	rtn := C.CreateBlob(s.ptr, C.int(functionCode))

	s.ptr = nil
	s.iso.ptr = nil

	raw_size := rtn.raw_size
	data := C.GoBytes(unsafe.Pointer(rtn.data), raw_size)

	C.SnapshotBlobDelete(rtn)

	return &StartupData{data: data, raw_size: raw_size}, nil
}

// Dispose deletes the reference to the SnapshotCreator.
func (s *SnapshotCreator) Dispose() {
	if s.ptr != nil {
		C.DeleteSnapshotCreator(s.ptr)
	}
}
