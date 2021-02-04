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
		ckey := C.CString(key)
		defer C.free(unsafe.Pointer(ckey))
		C.ObjectSet(o.ptr, ckey, value.ptr)
		return nil
	}

	C.ObjectSetIdx(o.ptr, C.uint32_t(idx), value.ptr)
	return nil
}

// Get tries to get a Value for a given Object property key.
func (o *Object) Get(key string) (*Value, error) {
	if o.ctx == nil {
		return nil, errors.New("v8go: unable to set property: Object has no Context")
	}
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
	panic("not implemented")
}

// HasIdx returns true if the object has a value at the given index.
func (o *Object) HasIdx(idx uint32) bool {
	panic("not implemented")
}

// Delete returns true if successful in deleting a named property on the object.
func (o *Object) Delete(key string) bool {
	panic("not implemented")
}

// DeleteIdx returns true if successful in deleting a value at a given index of the object.
func (o *Object) DeleteIdx(idx uint32) bool {
	panic("not implemented")
}
