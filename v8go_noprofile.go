// Copyright 2023 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

//go:build !v8goprofile

package v8go

func addIsolate(iso *Isolate) {}
func delIsolate(iso *Isolate) {}

func addContext(ctx *Context) {}
func delContext(ctx *Context) {}
