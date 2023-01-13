// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	v8 "github.com/sundeck-io/v8go"
)

func TestUnboundScriptRun_OnlyInTheSameIsolate(t *testing.T) {
	str := "function foo() { return 'bar'; }; foo()"
	i1 := v8.NewIsolate()
	defer i1.Dispose()

	us, err := i1.CompileUnboundScript(str, "script.js", v8.CompileOptions{})
	fatalIf(t, err)

	c1 := v8.NewContext(i1)
	defer c1.Close()

	val, err := us.Run(c1)
	fatalIf(t, err)
	if val.String() != "bar" {
		t.Fatalf("invalid value returned, expected bar got %v", val)
	}

	c2 := v8.NewContext(i1)
	defer c2.Close()

	val, err = us.Run(c2)
	fatalIf(t, err)
	if val.String() != "bar" {
		t.Fatalf("invalid value returned, expected bar got %v", val)
	}

	i2 := v8.NewIsolate()
	defer i2.Dispose()
	i2c1 := v8.NewContext(i2)
	defer i2c1.Close()

	if recoverPanic(func() { us.Run(i2c1) }) == nil {
		t.Error("expected panic running unbound script in a context belonging to a different isolate")
	}
}
