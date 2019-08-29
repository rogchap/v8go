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

func (i *Isolate) release() {
	C.IsolateRelease(i.ptr)
	i.ptr = nil
	runtime.SetFinalizer(i, nil)
}

func NewIsolate() *Isolate {
	v8once.Do(func() {
		C.Init()
	})
	iso := &Isolate{C.NewIsolate()}
	runtime.SetFinalizer(iso, (*Isolate).release)
	return iso
}
