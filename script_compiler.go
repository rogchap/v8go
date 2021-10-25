// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

type ScriptCompilerCachedData []byte

type ScriptCompilerCompileOption int

const (
	ScriptCompilerCompileOptionNoCompileOptions = iota
	ScriptCompilerCompileOptionConsumeCodeCache
	ScriptCompilerCompileOptionEagerCompile
)
