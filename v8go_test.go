// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"regexp"
	"testing"

	"rogchap.com/v8go"
)

func TestVersion(t *testing.T) {
	t.Parallel()
	rgx := regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+-v8go$`)
	v := v8go.Version()
	if !rgx.MatchString(v) {
		t.Errorf("version string is in the incorrect format: %s", v)
	}
}
