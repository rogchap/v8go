package v8go

//go:generate clang-format -i --verbose -style=Chromium v8go.h v8go.cc

// #cgo CXXFLAGS: -fno-rtti -fpic -std=c++14 -DV8_COMPRESS_POINTERS -DV8_31BIT_SMIS_ON_64BIT_ARCH
// #cgo darwin linux CXXFLAGS: -I${SRCDIR}/deps/include
// #cgo LDFLAGS: -pthread -lv8
// #cgo windows LDFLAGS: -lv8_libplatform
// #cgo darwin LDFLAGS: -L${SRCDIR}/deps/darwin-x86_64
// #cgo linux LDFLAGS: -L${SRCDIR}/deps/linux-x86_64
import "C"
