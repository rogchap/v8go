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

const (
	FunctionCodeHandlingKlear FunctionCodeHandling = iota
	FunctionCodeHandlingKeep
)

type StartupData struct {
	data     []byte
	raw_size C.int
}

type SnapshotCreator struct {
	ptr                 C.SnapshotCreatorPtr
	iso                 *Isolate
	defaultContextAdded bool
}

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

func (s *SnapshotCreator) GetIsolate() (*Isolate, error) {
	if s.ptr == nil {
		return nil, errors.New("v8go: Cannot get Isolate after creating the blob")
	}

	return s.iso, nil
}

func (s *SnapshotCreator) SetDeafultContext(ctx *Context) error {
	if s.defaultContextAdded {
		return errors.New("v8go: Cannot set multiple default context for snapshot creator")
	}

	C.SetDefaultContext(s.ptr, ctx.ptr)
	s.defaultContextAdded = true
	ctx.ptr = nil

	return nil
}

func (s *SnapshotCreator) AddContext(ctx *Context) (int, error) {
	if s.ptr == nil {
		return 0, errors.New("v8go: Cannot add context to snapshot creator after creating the blob")
	}

	index := C.AddContext(s.ptr, ctx.ptr)
	ctx.ptr = nil

	return int(index), nil
}

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

func (s *SnapshotCreator) Dispose() {
	if s.ptr != nil {
		C.DeleteSnapshotCreator(s.ptr)
	}
}
