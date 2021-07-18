// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
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

// Object is a JavaScript object (ECMA-262, 4.3.3)
type Object struct {
	*Value
}

// Instantiate a new blank object without any properties.
// Add properties on the object using Set.
func NewObject(ctx *ExecContext) *Object {
	return &Object{&Value{C.NewObject(ctx.ptr), ctx}}
}

// Set will set a property on the Object to a given value.
func (o *Object) Set(key string, val Valuer) error {
	if len(key) == 0 {
		return errors.New("v8go: You must provide a valid property key")
	}
	value := val.value()
	if value == nil {
		return errors.New("empty valuer given")
	}
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	C.ObjectSet(o.ptr, ckey, value.ptr)
	return nil
}

// Set will set a given index on the Object to a given value.
func (o *Object) SetIdx(idx uint32, val Valuer) error {
	value := val.value()
	if value == nil {
		return errors.New("empty valuer given")
	}
	C.ObjectSetIdx(o.ptr, C.uint32_t(idx), value.ptr)
	return nil
}

// Get tries to get a Value for a given Object property key.
func (o *Object) Get(key string) (*Value, error) {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	rtn := C.ObjectGet(o.ptr, ckey)
	return getValue(o.ctx, rtn), getError(rtn)
}

// GetIdx tries to get a Value at a give Object index.
func (o *Object) GetIdx(idx uint32) (*Value, error) {
	rtn := C.ObjectGetIdx(o.ptr, C.uint32_t(idx))
	return getValue(o.ctx, rtn), getError(rtn)
}

// Has calls the abstract operation HasProperty(O, P) described in ECMA-262, 7.3.10.
// Returns true, if the object has the property, either own or on the prototype chain.
func (o *Object) Has(key string) bool {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	return C.ObjectHas(o.ptr, ckey) != 0
}

// HasIdx returns true if the object has a value at the given index.
func (o *Object) HasIdx(idx uint32) bool {
	return C.ObjectHasIdx(o.ptr, C.uint32_t(idx)) != 0
}

// Delete returns true if successful in deleting a named property on the object.
func (o *Object) Delete(key string) bool {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	return C.ObjectDelete(o.ptr, ckey) != 0
}

// DeleteIdx returns true if successful in deleting a value at a given index of the object.
func (o *Object) DeleteIdx(idx uint32) bool {
	return C.ObjectDeleteIdx(o.ptr, C.uint32_t(idx)) != 0
}
