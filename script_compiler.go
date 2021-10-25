// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include "v8go.h"
import "C"
import "unsafe"

type ScriptCompilerCompileOption int

const (
	ScriptCompilerCompileOptionNoCompileOptions = iota
	ScriptCompilerCompileOptionConsumeCodeCache
	ScriptCompilerCompileOptionEagerCompile
)

type ScriptCompilerCachedData struct {
	ptr *C.ScriptCompilerCachedData
}

func (s *ScriptCompilerCachedData) Bytes() []byte {
	return []byte(C.GoBytes(unsafe.Pointer(s.ptr.data), s.ptr.length))
}
