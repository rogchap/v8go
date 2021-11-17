// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

//go:generate clang-format -i --verbose -style=Chromium v8go.h v8go.cc

// #cgo CXXFLAGS: -fno-rtti -fpic -std=c++14 -DV8_COMPRESS_POINTERS -DV8_31BIT_SMIS_ON_64BIT_ARCH -I${SRCDIR}/deps/include -Wall
// #cgo LDFLAGS: -pthread -lv8
// #cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/deps/darwin_x86_64
// #cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/deps/darwin_arm64
// #cgo linux LDFLAGS: -L${SRCDIR}/deps/linux_x86_64
import "C"

// These imports forces `go mod vendor` to pull in all the folders that
// contain V8 libraries and headers which otherwise would be ignored.
// DO NOT REMOVE
import (
	_ "rogchap.com/v8go/deps/darwin_arm64"
	_ "rogchap.com/v8go/deps/darwin_x86_64"
	_ "rogchap.com/v8go/deps/include"
	_ "rogchap.com/v8go/deps/linux_x86_64"
)
