// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"errors"
)

// PromiseState is the state of the Promise.
type PromiseState int

const (
	Pending PromiseState = iota
	Fulfilled
	Rejected
)

// PromiseResolver is the resolver object for the promise.
// Most cases will create a new PromiseResolver and return
// the associated Promise from the resolver.
type PromiseResolver struct {
	*Object
	prom *Promise
}

// Promise is the JavaScript promise object defined in ES6
type Promise struct {
	*Object
}

// MewPromiseResolver creates a new Promise resolver for the given context.
// The associated Promise will be in a Pending state.
func NewPromiseResolver(ctx *Context) (*PromiseResolver, error) {
	if ctx == nil {
		return nil, errors.New("v8go: Context is required")
	}
	rtn := C.NewPromiseResolver(ctx.ptr)
	obj, err := objectResult(ctx, rtn)
	if err != nil {
		return nil, err
	}
	return &PromiseResolver{obj, nil}, nil
}

// GetPromise returns the associated Promise object for this resolver.
// The Promise object is unique to the resolver and returns the same object
// on multiple calls.
func (r *PromiseResolver) GetPromise() *Promise {
	if r.prom == nil {
		ptr := C.PromiseResolverGetPromise(r.ptr)
		val := &Value{ptr, r.ctx}
		r.prom = &Promise{&Object{val}}
	}
	return r.prom
}

// Resolve invokes the Promise resolve state with the given value.
// The Promise state will transition from Pending to Fulfilled.
func (r *PromiseResolver) Resolve(val Valuer) bool {
	return C.PromiseResolverResolve(r.ptr, val.value().ptr) != 0
}

// Reject invokes the Promise reject state with the given value.
// The Promise state will transition from Pending to Rejected.
func (r *PromiseResolver) Reject(err *Value) bool {
	return C.PromiseResolverReject(r.ptr, err.ptr) != 0
}

// State returns the current state of the Promise.
func (p *Promise) State() PromiseState {
	return PromiseState(C.PromiseState(p.ptr))
}

// Result is the value result of the Promise. The Promise must
// NOT be in a Pending state, otherwise may panic. Call promise.State()
// to validate state before calling for the result.
func (p *Promise) Result() *Value {
	ptr := C.PromiseResult(p.ptr)
	val := &Value{ptr, p.ctx}
	return val
}

// Then accepts 1 or 2 callbacks.
// The first is invoked when the promise has been fulfilled.
// The second is invoked when the promise has been rejected.
// The returned Promise resolves after the callback finishes execution.
//
// V8 only invokes the callback when processing "microtasks".
// The default MicrotaskPolicy processes them when the call depth decreases to 0.
// Call (*Context).PerformMicrotaskCheckpoint to trigger it manually.
func (p *Promise) Then(cbs ...FunctionCallback) *Promise {
	var rtn C.RtnValue
	switch len(cbs) {
	case 1:
		cbID := p.ctx.iso.registerCallback(cbs[0])
		rtn = C.PromiseThen(p.ptr, C.int(cbID))
	case 2:
		cbID1 := p.ctx.iso.registerCallback(cbs[0])
		cbID2 := p.ctx.iso.registerCallback(cbs[1])
		rtn = C.PromiseThen2(p.ptr, C.int(cbID1), C.int(cbID2))

	default:
		panic("1 or 2 callbacks required")
	}
	obj, err := objectResult(p.ctx, rtn)
	if err != nil {
		panic(err) // TODO: Return error
	}
	return &Promise{obj}
}

// Catch invokes the given function if the promise is rejected.
// See Then for other details.
func (p *Promise) Catch(cb FunctionCallback) *Promise {
	cbID := p.ctx.iso.registerCallback(cb)
	rtn := C.PromiseCatch(p.ptr, C.int(cbID))
	obj, err := objectResult(p.ctx, rtn)
	if err != nil {
		panic(err) // TODO: Return error
	}
	return &Promise{obj}
}
