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
	ptr := C.NewPromiseResolver(ctx.ptr)
	val := &Value{ptr, ctx}
	return &PromiseResolver{&Object{val}, nil}, nil
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
	r.ctx.register()
	defer r.ctx.deregister()
	return C.PromiseResolverResolve(r.ptr, val.value().ptr) != 0
}

// Reject invokes the Promise reject state with the given value.
// The Promise state will transition from Pending to Rejected.
func (r *PromiseResolver) Reject(err *Value) bool {
	r.ctx.register()
	defer r.ctx.deregister()
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

// Then invokes the given function when the promise has been fulfilled, not rejected.
// The returned Promise resolves after the given function has finished execution.
func (p *Promise) Then(cb FunctionCallback) *Promise {
	p.ctx.register()
	defer p.ctx.deregister()
	cbID := p.ctx.iso.registerCallback(cb)
	ptr := C.PromiseThen(p.ptr, C.int(cbID))
	return &Promise{&Object{&Value{ptr, p.ctx}}}
}

// Then2 invokes one of the given functions when the promise is fulfilled or rejected.
// The returned Promise resolves after the callback has finished execution.
func (p *Promise) Then2(onFulfilled, onRejected FunctionCallback) *Promise {
	p.ctx.register()
	defer p.ctx.deregister()
	onFulfilledID := p.ctx.iso.registerCallback(onFulfilled)
	onRejectedID := p.ctx.iso.registerCallback(onRejected)
	ptr := C.PromiseThen2(p.ptr, C.int(onFulfilledID), C.int(onRejectedID))
	return &Promise{&Object{&Value{ptr, p.ctx}}}
}

// Catch invokes the given function if the promise is rejected.
// The returned Promise resolves after the callback has finished execution.
func (p *Promise) Catch(cb FunctionCallback) *Promise {
	p.ctx.register()
	defer p.ctx.deregister()
	cbID := p.ctx.iso.registerCallback(cb)
	ptr := C.PromiseCatch(p.ptr, C.int(cbID))
	return &Promise{&Object{&Value{ptr, p.ctx}}}
}
