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

// PropertyAttribute are the attribute flags for a property on an Object.
// Typical usage when setting an Object or TemplateObject property, and
// can also be validated when accessing a property.
type PropertyAttribute uint8

const (
	// None.
	None PropertyAttribute = 0
	// ReadOnly, ie. not writable.
	ReadOnly PropertyAttribute = 1 << iota
	// DontEnum, ie. not enumerable.
	DontEnum
	// DontDelete, ie. not configurable.
	DontDelete
)

// ObjectTemplate is used to create objects at runtime.
// Properties added to an ObjectTemplate are added to each object created from the ObjectTemplate.
type ObjectTemplate struct {
	ptr C.ObjectTemplatePtr
	iso *Isolate
}

// NewObjectTemplate creates a new ObjectTemplate.
// The *ObjectTemplate can be used as a v8go.ContextOption to create a global object in a Context.
func NewObjectTemplate(iso *Isolate) (*ObjectTemplate, error) {
	if iso == nil {
		return nil, errors.New("v8go: failed to create new ObjectTemplate: Isolate cannot be <nil>")
	}
	ob := &ObjectTemplate{
		ptr: C.NewObjectTemplate(iso.ptr),
		iso: iso,
	}
	runtime.SetFinalizer(ob, (*ObjectTemplate).finalizer)
	return ob, nil
}

// Set adds a property to each instance created by this template.
// The property must be defined either as a primitive value, or a template.
// If the value passed is a Go supported primitive (string, int32, uint32, int64, uint64, float64, big.Int)
// then a value will be created and set as the value property.
func (o *ObjectTemplate) Set(name string, val interface{}, attributes ...PropertyAttribute) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var attrs PropertyAttribute
	for _, a := range attributes {
		attrs |= a
	}

	switch v := val.(type) {
	case string, int32, uint32, int64, uint64, float64, bool, *big.Int:
		newVal, err := NewValue(o.iso, v)
		if err != nil {
			return fmt.Errorf("v8go: unable to create new value: %v", err)
		}
		C.ObjectTemplateSetValue(o.ptr, cname, newVal.ptr, C.int(attrs))
	case *ObjectTemplate:
		C.ObjectTemplateSetObjectTemplate(o.ptr, cname, v.ptr, C.int(attrs))
	case *Value:
		if v.IsObject() || v.IsExternal() {
			return errors.New("v8go: unsupported property: value type must be a primitive or use a template")
		}
		C.ObjectTemplateSetValue(o.ptr, cname, v.ptr, C.int(attrs))
	default:
		return fmt.Errorf("v8go: unsupported property type `%T`, must be one of string, int32, uint32, int64, uint64, float64, *big.Int, *v8go.Value or *v8go.ObjectTemplate", v)
	}

	return nil
}

func (o *ObjectTemplate) apply(opts *contextOptions) {
	opts.gTmpl = o
}

func (o *ObjectTemplate) finalizer() {
	C.ObjectTemplateDispose(o.ptr)
	o.ptr = nil
	runtime.SetFinalizer(o, nil)
}
