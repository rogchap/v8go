// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

//go:build leakcheck
// +build leakcheck

package v8go_test

import (
	"os"
	"testing"

	"github.com/sundeck-io/v8go"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	v8go.DoLeakSanitizerCheck()
	os.Exit(exitCode)
}
