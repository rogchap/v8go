/*
Package v8go is an API wrapper to the v8 Javascript engine
*/
package v8go

// #include "v8go.h"
import "C"

func Version() string {
	return C.GoString(C.Version())
}
