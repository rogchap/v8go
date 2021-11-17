// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include "v8go.h"
import "C"

type CompileMode C.int

var (
	CompileModeDefault = CompileMode(C.ScriptCompilerNoCompileOptions)
	CompileModeEager   = CompileMode(C.ScriptCompilerEagerCompile)
)

type CompilerCachedData struct {
	Bytes    []byte
	Rejected bool
}
