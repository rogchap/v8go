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

// PropertyAttribute are the attributes for a property
type PropertyAttribute uint8

const (
	// None.
	None PropertyAttribute = 0
	// ReadOnly, ie. not writable.
	ReadOnly = 1 << iota
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
		return nil, errors.New("object_template: failed to create new ObjectTemplate: isolate cannot be <nil>")
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
// then a value will be created an set as the value property.
func (o *ObjectTemplate) Set(name string, val interface{}, attributes PropertyAttribute) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	switch v := val.(type) {
	case string, int32, uint32, int64, uint64, float64, bool, *big.Int:
		newVal, err := NewValue(o.iso, v)
		if err != nil {
			return fmt.Errorf("object_template: unable to create new value: %v", err)
		}
		C.ObjectTemplateSetValue(o.ptr, cname, newVal.ptr, C.int(attributes))
	case *ObjectTemplate:
		C.ObjectTemplateSetObjectTemplate(o.ptr, cname, v.ptr, C.int(attributes))
	case *Value:
		if v.IsObject() || v.IsExternal() {
			return errors.New("object_template: unsupported property: value type must be a primitive or use a template")
		}
		C.ObjectTemplateSetValue(o.ptr, cname, v.ptr, C.int(attributes))
	default:
		return fmt.Errorf("object_template: unsupported property type `%T`, must be one of string, int32, uint32, int64, uint64, float64, *big.Int, *v8go.Value or *v8go.ObjectTemplate", v)
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
