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
	ptr C.ObjectTemplatePtr
	iso *Isolate
}

func newTemplate(iso *Isolate) (*template, error) {
	if iso == nil {
		return nil, errors.New("v8go: failed to create new Template: Isolate cannot be <nil>")
	}

	ob := &template{
		ptr: C.NewObjectTemplate(iso.ptr),
		iso: iso,
	}
	runtime.SetFinalizer(ob, (*template).finalizer)
	return ob, nil
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
		C.ObjectTemplateSetValue(t.ptr, cname, newVal.ptr, C.int(attrs))
	case *ObjectTemplate:
		C.ObjectTemplateSetObjectTemplate(t.ptr, cname, v.ptr, C.int(attrs))
	case *Value:
		if v.IsObject() || v.IsExternal() {
			return errors.New("v8go: unsupported property: value type must be a primitive or use a template")
		}
		C.ObjectTemplateSetValue(t.ptr, cname, v.ptr, C.int(attrs))
	default:
		return fmt.Errorf("v8go: unsupported property type `%T`, must be one of string, int32, uint32, int64, uint64, float64, *big.Int, *v8go.Value, *v8go.ObjectTemplate or *v8go.FunctionTemplate", v)
	}

	return nil
}

func (o *template) finalizer() {
	C.ObjectTemplateDispose(o.ptr)
	o.ptr = nil
	runtime.SetFinalizer(o, nil)
}
