// Copyright 2023 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

//go:build v8goprofile

package v8go

import "runtime/pprof"

func getProfile(profileName string) *pprof.Profile {
	if p := pprof.Lookup(); p != nil {
		return p
	}

	return pprof.NewProfile(profileName)
}

func addIsolate(iso *Isolate) {
	getProfile("rogchap.com/v8go/gv8go.Isolate").Add(iso, 1)
}

func delIsolate(iso *Isolate) {
	getProfile("rogchap.com/v8go/gv8go.Isolate").Remove(iso)
}

func addContext(ctx *Context) {
	getProfile("rogchap.com/v8go/gv8go.Context").Add(ctx, 1)
}

func delContext(ctx *Context) {
	getProfile("rogchap.com/v8go/gv8go.Context").Remove(ctx)
}
