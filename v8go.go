/*
Package v8go provides an API to execute JavaScript.
*/
package v8go

// #include "v8go.h"
import "C"

// Version returns the version of the V8 Engine with the -v8go suffix
func Version() string {
	return C.GoString(C.Version())
}
