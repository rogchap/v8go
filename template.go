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
	"runtime"
	"unsafe"
)

type template struct {
	ptr C.TemplatePtr
	iso *Isolate
}

// Set adds a property to each instance created by this template.
// The property must be defined either as a primitive value, or a template.
// If the value passed is a Go supported primitive (string, int32, uint32, int64, uint64, float64, big.Int)
// then a value will be created and set as the value property.
func (t *template) Set(name string, val interface{}, attributes ...PropertyAttribute) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var attrs PropertyAttribute
	for _, a := range attributes {
		attrs |= a
	}

	switch v := val.(type) {
	case string, int32, uint32, int64, uint64, float64, bool, *big.Int:
		newVal, err := NewValue(t.iso, v)
		if err != nil {
			return fmt.Errorf("v8go: unable to create new value: %v", err)
		}
		C.TemplateSetValue(t.ptr, cname, newVal.ptr, C.int(attrs))
	case *ObjectTemplate:
		C.TemplateSetTemplate(t.ptr, cname, v.ptr, C.int(attrs))
		runtime.KeepAlive(v)
	case *FunctionTemplate:
		C.TemplateSetTemplate(t.ptr, cname, v.ptr, C.int(attrs))
		runtime.KeepAlive(v)
	case *Value:
		if v.IsObject() || v.IsExternal() {
			return errors.New("v8go: unsupported property: value type must be a primitive or use a template")
		}
		C.TemplateSetValue(t.ptr, cname, v.ptr, C.int(attrs))
	default:
		return fmt.Errorf("v8go: unsupported property type `%T`, must be one of string, int32, uint32, int64, uint64, float64, *big.Int, *v8go.Value, *v8go.ObjectTemplate or *v8go.FunctionTemplate", v)
	}
	runtime.KeepAlive(t)

	return nil
}

// SetSymbol adds a property to each instance created by this template.
// The property must be defined either as a primitive value, or a template.
// If the value passed is a Go supported primitive (string, int32, uint32, int64, uint64, float64, big.Int)
// then a value will be created and set as the value property.
func (t *template) SetSymbol(key *Symbol, val interface{}, attributes ...PropertyAttribute) error {
	var attrs PropertyAttribute
	for _, a := range attributes {
		attrs |= a
	}

	switch v := val.(type) {
	case string, int32, uint32, int64, uint64, float64, bool, *big.Int:
		newVal, err := NewValue(t.iso, v)
		if err != nil {
			return fmt.Errorf("v8go: unable to create new value: %v", err)
		}
		if C.TemplateSetAnyValue(t.ptr, key.ptr, newVal.ptr, C.int(attrs)) == 0 {
			return fmt.Errorf("v8go: unable to set property for symbol %v", key)
		}
	case *ObjectTemplate:
		if C.TemplateSetAnyTemplate(t.ptr, key.ptr, v.ptr, C.int(attrs)) == 0 {
			return fmt.Errorf("v8go: unable to set property for symbol %v", key)
		}
		runtime.KeepAlive(v)
	case *FunctionTemplate:
		if C.TemplateSetAnyTemplate(t.ptr, key.ptr, v.ptr, C.int(attrs)) == 0 {
			return fmt.Errorf("v8go: unable to set property for symbol %v", key)
		}
		runtime.KeepAlive(v)
	case *Value:
		if v.IsObject() || v.IsExternal() {
			return errors.New("v8go: unsupported property: value type must be a primitive or use a template")
		}
		if C.TemplateSetAnyValue(t.ptr, key.ptr, v.ptr, C.int(attrs)) == 0 {
			return fmt.Errorf("v8go: unable to set property for symbol %v", key)
		}
	default:
		return fmt.Errorf("v8go: unsupported property type `%T`, must be one of string, int32, uint32, int64, uint64, float64, *big.Int, *v8go.Value, *v8go.ObjectTemplate or *v8go.FunctionTemplate", v)
	}
	runtime.KeepAlive(t)

	return nil
}

func (t *template) finalizer() {
	// Using v8::PersistentBase::Reset() wouldn't be thread-safe to do from
	// this finalizer goroutine so just free the wrapper and let the template
	// itself get cleaned up when the isolate is disposed.
	C.TemplateFreeWrapper(t.ptr)
	t.ptr = nil
}
