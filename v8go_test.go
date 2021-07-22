// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"regexp"
	"testing"

	"github.com/Shopify/v8go"
)

func TestVersion(t *testing.T) {
	t.Parallel()
	rgx := regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+-v8go$`)
	v := v8go.Version()
	if !rgx.MatchString(v) {
		t.Errorf("version string is in the incorrect format: %s", v)
	}
}

func TestSetFlag(t *testing.T) {
	t.Parallel()
	ctx, _ := v8go.NewContext()
	if _, err := ctx.RunScript("a = 1", "default.js"); err != nil {
		t.Errorf("expected <nil> error, but got: %v", err)
	}
	v8go.SetFlags("--use_strict")
	if _, err := ctx.RunScript("b = 1", "use_strict.js"); err == nil {
		t.Error("expected error but got <nil>")
	}
	v8go.SetFlags("--nouse_strict")
	if _, err := ctx.RunScript("c = 1", "nouse_strict.js"); err != nil {
		t.Errorf("expected <nil> error, but got: %v", err)
	}
}
