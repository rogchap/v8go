// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

//go:build !muslgcc || glibc
// +build !muslgcc glibc

package v8go

// #cgo linux,amd64 LDFLAGS: -L${SRCDIR}/deps/linux_x86_64 -ldl
import "C"
