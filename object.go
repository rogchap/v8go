// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"errors"
	"fmt"
	"math/big"
	"unsafe"
)

// Object is a JavaScript object (ECMA-262, 4.3.3)
type Object struct {
	*Value
}

func (o *Object) MethodCall(methodName string, args ...Valuer) (*Value, error) {
	ckey := C.CString(methodName)
	defer C.free(unsafe.Pointer(ckey))

	getRtn := C.ObjectGet(o.ptr, ckey)
	prop, err := valueResult(o.ctx, getRtn)
	if err != nil {
		return nil, err
	}
	fn, err := prop.AsFunction()
	if err != nil {
		return nil, err
	}
	return fn.Call(o, args...)
}

func coerceValue(iso *Isolate, val interface{}) (*Value, error) {
	switch v := val.(type) {
	case string, int32, uint32, int64, uint64, float64, bool, *big.Int:
		// ignoring error as code cannot reach the error state as we are already
		// validating the new value types in this case statement
		value, _ := NewValue(iso, v)
		return value, nil
	case Valuer:
		return v.value(), nil
	default:
		return nil, fmt.Errorf("v8go: unsupported object property type `%T`", v)
	}
}

// Set will set a property on the Object to a given value.
// Supports all value types, eg: Object, Array, Date, Set, Map etc
// If the value passed is a Go supported primitive (string, int32, uint32, int64, uint64, float64, big.Int)
// then a *Value will be created and set as the value property.
func (o *Object) Set(key string, val interface{}) error {
	if len(key) == 0 {
		return errors.New("v8go: You must provide a valid property key")
	}

	value, err := coerceValue(o.ctx.iso, val)
	if err != nil {
		return err
	}

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	C.ObjectSet(o.ptr, ckey, value.ptr)
	return nil
}

// SetSymbol will set a property on the Object to a given value.
// Supports all value types, eg: Object, Array, Date, Set, Map etc
// If the value passed is a Go supported primitive (string, int32, uint32, int64, uint64, float64, big.Int)
// then a *Value will be created and set as the value property.
func (o *Object) SetSymbol(key *Symbol, val interface{}) error {
	value, err := coerceValue(o.ctx.iso, val)
	if err != nil {
		return err
	}

	C.ObjectSetAnyKey(o.ptr, key.ptr, value.ptr)
	return nil
}

// Set will set a given index on the Object to a given value.
// Supports all value types, eg: Object, Array, Date, Set, Map etc
// If the value passed is a Go supported primitive (string, int32, uint32, int64, uint64, float64, big.Int)
// then a *Value will be created and set as the value property.
func (o *Object) SetIdx(idx uint32, val interface{}) error {
	value, err := coerceValue(o.ctx.iso, val)
	if err != nil {
		return err
	}

	C.ObjectSetIdx(o.ptr, C.uint32_t(idx), value.ptr)
	return nil
}

// Get tries to get a Value for a given Object property key.
func (o *Object) Get(key string) (*Value, error) {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	rtn := C.ObjectGet(o.ptr, ckey)
	return valueResult(o.ctx, rtn)
}

// GetSymbol tries to get a Value for a given Object property key.
func (o *Object) GetSymbol(key *Symbol) (*Value, error) {
	rtn := C.ObjectGetAnyKey(o.ptr, key.ptr)
	return valueResult(o.ctx, rtn)
}

// GetIdx tries to get a Value at a give Object index.
func (o *Object) GetIdx(idx uint32) (*Value, error) {
	rtn := C.ObjectGetIdx(o.ptr, C.uint32_t(idx))
	return valueResult(o.ctx, rtn)
}

// Has calls the abstract operation HasProperty(O, P) described in ECMA-262, 7.3.10.
// Returns true, if the object has the property, either own or on the prototype chain.
func (o *Object) Has(key string) bool {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	return C.ObjectHas(o.ptr, ckey) != 0
}

// HasSymbol calls the abstract operation HasProperty(O, P) described in ECMA-262, 7.3.10.
// Returns true, if the object has the property, either own or on the prototype chain.
func (o *Object) HasSymbol(key *Symbol) bool {
	return C.ObjectHasAnyKey(o.ptr, key.ptr) != 0
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

// DeleteSymbol returns true if successful in deleting a named property on the object.
func (o *Object) DeleteSymbol(key *Symbol) bool {
	return C.ObjectDeleteAnyKey(o.ptr, key.ptr) != 0
}

// DeleteIdx returns true if successful in deleting a value at a given index of the object.
func (o *Object) DeleteIdx(idx uint32) bool {
	return C.ObjectDeleteIdx(o.ptr, C.uint32_t(idx)) != 0
}
