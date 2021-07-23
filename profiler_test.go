// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	"github.com/Shopify/v8go"
)

func TestProfiler(t *testing.T) {
	// t.Parallel()

	ctx, _ := v8go.NewContext()
	profiler, _ := v8go.NewProfiler(ctx)

	profiler.Start()
	ctx.RunScript("const foo = {}; foo", "")
	jsonProfile, err := profiler.Stop()
	if err != nil {
		t.Errorf("Invalid json profile: %v", err)
	}
	if len(jsonProfile) == 0 {
		t.Error("Missing profile data")
	}
}
