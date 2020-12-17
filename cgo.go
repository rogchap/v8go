package v8go

// #cgo CXXFLAGS: -fno-rtti -fpic -std=c++14 -I${SRCDIR}/deps/include -DV8_COMPRESS_POINTERS -DV8_31BIT_SMIS_ON_64BIT_ARCH
// #cgo LDFLAGS: -pthread -lv8
// #cgo darwin LDFLAGS: -L${SRCDIR}/deps/darwin-x86_64
// #cgo linux LDFLAGS: -L${SRCDIR}/deps/linux-x86_64
import "C"
