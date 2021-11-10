// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

//go:build leakcheck
// +build leakcheck

package v8go

// #cgo CPPFLAGS: -fsanitize=address
// #cgo LDFLAGS: -fsanitize=address
//
// #include <sanitizer/lsan_interface.h>
import "C"

import "runtime"

// Call LLVM Leak Sanitizer's at-exit hook that doesn't
// get called automatically by Go.
func DoLeakSanitizerCheck() {
	runtime.GC()
	C.__lsan_do_leak_check()
}
