package v8go

// #include "v8go.h"
import "C"

import (
	"runtime"
	"sync"
)

var v8once sync.Once

type Isolate struct {
	ptr C.IsolatePtr
}

func NewIsolate() *Isolate {
	v8once.Do(func() {
		C.Init()
	})
	iso := &Isolate{C.NewIsolate()}
	runtime.SetFinalizer(iso, (*Isolate).finalizer)
	return iso
}

func (i *Isolate) finalizer() {
	C.IsolateDispose(i.ptr)
	i.ptr = nil
	runtime.SetFinalizer(i, nil)
}
