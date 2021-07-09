// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

// Templater allows composition of templates and will
// return actual *Value if context is provided.
type Templater interface {
	ContextValue(ctx *ExecContext) (Valuer, error)
	Set(name string, val Templater, attributes ...PropertyAttribute) error
}

type template struct {
	ptr C.TemplatePtr
	iso *Isolate
}

// Set adds a property to each instance created by this template.
// The property must be defined either as a primitive value, or a template.
// If the value passed is a Go supported primitive (string, int32, uint32, int64, uint64, float64, big.Int)
// then a value will be created and set as the value property.
func (t *template) Set(name string, val Templater, attributes ...PropertyAttribute) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var attrs PropertyAttribute
	for _, a := range attributes {
		attrs |= a
	}

	switch v := val.(type) {
	case *ObjectTemplate:
		C.TemplateSetTemplate(t.ptr, cname, v.ptr, C.int(attrs))
	case *FunctionTemplate:
		C.TemplateSetTemplate(t.ptr, cname, v.ptr, C.int(attrs))
	case *ValueTemplate:
		C.TemplateSetValue(t.ptr, cname, v.Value.ptr, C.int(attrs))
	default:
		return fmt.Errorf("v8go: unsupported property type `%T`, must be one of string, int32, uint32, int64, uint64, float64, *big.Int, *v8go.Value, *v8go.ObjectTemplate or *v8go.FunctionTemplate", v)
	}

	return nil
}

func (t *template) finalizer() {
	C.TemplateFree(t.ptr)
	t.ptr = nil
	runtime.SetFinalizer(t, nil)
}
