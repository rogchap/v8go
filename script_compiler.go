// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include "v8go.h"
import "C"

type ScriptCompilerCompileOption C.int

var (
	ScriptCompilerNoCompileOptions = ScriptCompilerCompileOption(C.ScriptCompilerNoCompileOptions)
	ScriptCompilerEagerCompile = ScriptCompilerCompileOption(C.ScriptCompilerEagerCompile)
)

type ScriptCompilerCachedData struct {
	Bytes    []byte
	Rejected bool
}
