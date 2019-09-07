package v8go

// #include "v8go.h"
import "C"

import (
	"runtime"
	"sync"
)

var v8once sync.Once

// Isolate is a JavaScript VM instance with its own heap and
// garbage collector. Most applications will create one isolate
// with many V8 contexts for execution.
type Isolate struct {
	ptr C.IsolatePtr
}

// NewIsolate creates a new V8 isolate. Only one thread may access
// a given isolate at a time, but different threads may access
// different isolates simultaneously.
func NewIsolate() (*Isolate, error) {
	v8once.Do(func() {
		C.Init()
	})
	iso := &Isolate{C.NewIsolate()}
	runtime.SetFinalizer(iso, (*Isolate).finalizer)
	// TODO: [RC] catch any C++ exceptions and return as error
	return iso, nil
}

// Forcefully terminate the current thread of JavaScript execution
// in the given isolate.
func (i *Isolate) TerminateExecution() {
	C.IsolateTerminateExecution(i.ptr)
}

func (i *Isolate) finalizer() {
	C.IsolateDispose(i.ptr)
	i.ptr = nil
	runtime.SetFinalizer(i, nil)
}
