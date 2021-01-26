package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

type FunctionCallback func(info *FunctionCallbackInfo) *Value

type FunctionCallbackInfo struct {
	args []*Value
}

func (i *FunctionCallbackInfo) Args() []*Value {
	return i.args
}

type FunctionTemplate struct {
	*template
}

func NewFunctionTemplate(iso *Isolate, callback FunctionCallback) (*FunctionTemplate, error) {
	if iso == nil {
		return nil, errors.New("v8go: failed to create new FunctionTemplate: Isolate cannot be <nil>")
	}

	tmpl := &template{
		ptr: C.NewFunctionTemplate(iso.ptr, unsafe.Pointer(&callback)),
		iso: iso,
	}
	runtime.SetFinalizer(tmpl, (*template).finalizer)
	return &FunctionTemplate{tmpl}, nil
}

//export goFunctionCallback
func goFunctionCallback(callback unsafe.Pointer, args *C.ValuePtr, args_count int) {
	callbackFunc := *(*FunctionCallback)(callback)

	//TODO: This will need access to the Context to be able to create return values
	info := &FunctionCallbackInfo{
		args: make([]*Value, args_count),
	}

	argv := (*[1 << 30]C.ValuePtr)(unsafe.Pointer(args))[:args_count:args_count]
	for i, v := range argv {
		//TODO(rogchap): We must pass the current context and add this to the value struct
		val := &Value{ptr: v}
		runtime.SetFinalizer(val, (*Value).finalizer)
		info.args[i] = val
	}

	rtnVal := callbackFunc(info)
	//TODO: deal with the rtnVal
	_ = rtnVal
}
