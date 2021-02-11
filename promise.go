// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"errors"
	"runtime"
)

// PromiseState
type PromiseState int

const (
	Pending PromiseState = iota
	Fulfilled
	Rejected
)

// PromiseResolver
type PromiseResolver struct {
	*Object
	prom *Promise
}

// Promise
type Promise struct {
	*Object
}

// MewPromiseResolver
func NewPromiseResolver(ctx *Context) (*PromiseResolver, error) {
	if ctx == nil {
		return nil, errors.New("v8go: Context is required")
	}
	ptr := C.NewPromiseResolver(ctx.ptr)
	val := &Value{ptr, ctx}
	runtime.SetFinalizer(val, (*Value).finalizer)
	return &PromiseResolver{&Object{val}, nil}, nil
}

// GetPromise
func (r *PromiseResolver) GetPromise() *Promise {
	if r.prom == nil {
		ptr := C.PromiseResolverGetPromise(r.ptr)
		val := &Value{ptr, r.ctx}
		runtime.SetFinalizer(val, (*Value).finalizer)
		r.prom = &Promise{&Object{val}}
	}
	return r.prom
}

// Resolve
func (r *PromiseResolver) Resolve(val Valuer) bool {
	r.ctx.register()
	defer r.ctx.deregister()
	return C.PromiseResolverResolve(r.ptr, val.value().ptr) != 0
}

// Reject
func (r *PromiseResolver) Reject(err *Value) bool {
	r.ctx.register()
	defer r.ctx.deregister()
	return C.PromiseResolverReject(r.ptr, err.ptr) != 0
}

// State
func (p *Promise) State() PromiseState {
	return PromiseState(C.PromiseState(p.ptr))
}

// Result
func (p *Promise) Result() *Value {
	ptr := C.PromiseResult(p.ptr)
	val := &Value{ptr, p.ctx}
	runtime.SetFinalizer(val, (*Value).finalizer)
	return val
}
