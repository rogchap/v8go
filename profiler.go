// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"

import (
	"encoding/json"
	"errors"
	"runtime"
	"unsafe"
)

type ProfileResult struct {
	Profile json.RawMessage `json:"profile"`
}

// TODO: Handle error messages
type ProfileMessage struct {
	// ignore id field
	Result ProfileResult `json:"result"`
}

func extractProfileData(jsonInspectorMessage string) (string, error) {
	var profile ProfileMessage
	err := json.Unmarshal([]byte(jsonInspectorMessage), &profile)
	if err != nil {
		return "", err
	} else {
		return string(profile.Result.Profile), nil
	}
}

type Profiler struct {
	ptr C.ProfilerPtr
	ctx *Context
}

func NewProfiler(ctx *Context) (*Profiler, error) {
	if ctx == nil {
		return nil, errors.New("v8go: failed to create new Profiler: Context cannot be <nil>")
	}

	profiler := &Profiler{
		ptr: C.NewProfiler(ctx.ptr),
		ctx: ctx,
	}
	runtime.SetFinalizer(profiler, (*Profiler).finalizer)

	return profiler, nil
}

func (p *Profiler) Start() {
	C.ProfilerStart(p.ptr)
}

func (p *Profiler) Stop() (string, error) {
	var length C.int = 0
	s := C.ProfilerStop(p.ptr, &length)
	data, err := extractProfileData(C.GoStringN(s, length))
	defer C.free(unsafe.Pointer(s))
	return data, err
}

func (p *Profiler) finalizer() {
	C.ProfilerFree(p.ptr)
	p.ptr = nil
	runtime.SetFinalizer(p, nil)
}
