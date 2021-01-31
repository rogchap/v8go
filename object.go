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

type Object struct {
	*Value
}

// Set will set a property on the Object to a given value.
// Supports all value types, eg: Object, Array, Date, Set, Map etc
// If the value passed is a Go supported primitive (string, int32, uint32, int64, uint64, float64, big.Int)
// then a *Value will be created and set as the value property.
func (o *Object) Set(key string, val interface{}) error {
	if o.ctx == nil {
		return errors.New("v8go: unable to set property: Object has no Context")
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var value *Value
	switch v := val.(type) {
	case string, int32, uint32, int64, uint64, float64, bool, *big.Int:
		var err error
		if value, err = NewValue(o.ctx.iso, v); err != nil {
			return fmt.Errorf("v8go: unable to create new value: %v", err)
		}
	case valuer:
		value = v.value()
	default:
		return fmt.Errorf("v8go: unsupported object property type `%T`", v)
	}

	C.ObjectSet(o.ptr, cname, value)

	return nil
}
