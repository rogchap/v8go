package v8go

// #include <stdlib.h>
//
// #cgo CXXFLAGS: -fno-rtti -fpic -std=c++14 -I${SRCDIR}/deps/include
// #cgo LDFLAGS: -pthread -lv8
// #cgo darwin LDFLAGS: -L${SRCDIR}/deps/darwin-x86_64
import "C"
