// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"errors"
	"strings"
	"testing"

	v8 "rogchap.com/v8go"
)

func TestNewError(t *testing.T) {
	t.Parallel()

	tsts := []struct {
		New      func(*v8.Isolate, string) *v8.Exception
		WantType string
	}{
		{v8.NewRangeError, "RangeError"},
		{v8.NewReferenceError, "ReferenceError"},
		{v8.NewSyntaxError, "SyntaxError"},
		{v8.NewTypeError, "TypeError"},
		{v8.NewWasmCompileError, "CompileError"},
		{v8.NewWasmLinkError, "LinkError"},
		{v8.NewWasmRuntimeError, "RuntimeError"},
		{v8.NewError, "Error"},
	}
	for _, tst := range tsts {
		t.Run(tst.WantType, func(t *testing.T) {
			iso := v8.NewIsolate()
			defer iso.Dispose()

			got := tst.New(iso, "amessage")
			if !got.IsNativeError() {
				t.Error("IsNativeError returned false, want true")
			}
			if got := got.Error(); !strings.Contains(got, " "+tst.WantType+":") {
				t.Errorf("Error(): got %q, want containing %q", got, tst.WantType)
			}
		})
	}
}

func TestExceptionAs(t *testing.T) {
	iso := v8.NewIsolate()
	defer iso.Dispose()

	want := v8.NewRangeError(iso, "faked error")

	var got *v8.Exception
	if !want.As(&got) {
		t.Fatalf("As failed")
	}

	if got != want {
		t.Errorf("As: got %#v, want %#v", got, want)
	}
}

func TestExceptionIs(t *testing.T) {
	iso := v8.NewIsolate()
	defer iso.Dispose()

	t.Run("ok", func(t *testing.T) {
		ex := v8.NewRangeError(iso, "faked error")
		if !ex.Is(v8.NewRangeError(iso, "faked error")) {
			t.Fatalf("Is: got false, want true")
		}
	})

	t.Run("notok", func(t *testing.T) {
		if (&v8.Exception{}).Is(errors.New("other error")) {
			t.Fatalf("Is: got true, want false")
		}
	})
}
