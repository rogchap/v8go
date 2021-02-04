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

// Set will set a property on the Object to a given value.
// Supports all value types, eg: Object, Array, Date, Set, Map etc
// If the value passed is a Go supported primitive (string, int32, uint32, int64, uint64, float64, big.Int)
// then a *Value will be created and set as the value property.
func (o *Object) Set(key string, val interface{}) error {
	if len(key) == 0 {
		return errors.New("v8go: You must provide a valid property key")
	}
	return set(o, key, 0, val)
}

// Set will set a given index on the Object to a given value.
// Supports all value types, eg: Object, Array, Date, Set, Map etc
// If the value passed is a Go supported primitive (string, int32, uint32, int64, uint64, float64, big.Int)
// then a *Value will be created and set as the value property.
func (o *Object) SetIdx(idx uint32, val interface{}) error {
	return set(o, "", idx, val)
}

func set(o *Object, key string, idx uint32, val interface{}) error {
	if o.ctx == nil {
		return errors.New("v8go: unable to set property: Object has no Context")
	}

	var value *Value
	switch v := val.(type) {
	case string, int32, uint32, int64, uint64, float64, bool, *big.Int:
		var err error
		if value, err = NewValue(o.ctx.iso, v); err != nil {
			return fmt.Errorf("v8go: unable to create new value: %v", err)
		}
	case Valuer:
		value = v.value()
	default:
		return fmt.Errorf("v8go: unsupported object property type `%T`", v)
	}

	if len(key) > 0 {
		cname := C.CString(key)
		defer C.free(unsafe.Pointer(cname))
		C.ObjectSet(o.ptr, cname, value.ptr)
		return nil
	}

	C.ObjectSetIdx(o.ptr, C.uint32_t(idx), value.ptr)
	return nil
}
