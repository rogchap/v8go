// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"fmt"
	"testing"

	v8 "rogchap.com/v8go"
)

func TestJSErrorFormat(t *testing.T) {
	t.Parallel()
	tests := [...]struct {
		name            string
		err             error
		defaultVerb     string
		defaultVerbFlag string
		stringVerb      string
		quoteVerb       string
	}{
		{"WithStack", &v8.JSError{Message: "msg", StackTrace: "stack"}, "msg", "stack", "msg", `"msg"`},
		{"WithoutStack", &v8.JSError{Message: "msg"}, "msg", "msg", "msg", `"msg"`},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if s := fmt.Sprintf("%v", tt.err); s != tt.defaultVerb {
				t.Errorf("incorrect format for %%v: %s", s)
			}
			if s := fmt.Sprintf("%+v", tt.err); s != tt.defaultVerbFlag {
				t.Errorf("incorrect format for %%+v: %s", s)
			}
			if s := fmt.Sprintf("%s", tt.err); s != tt.stringVerb {
				t.Errorf("incorrect format for %%s: %s", s)
			}
			if s := fmt.Sprintf("%q", tt.err); s != tt.quoteVerb {
				t.Errorf("incorrect format for %%q: %s", s)
			}
		})
	}
}

func TestJSErrorOutput(t *testing.T) {
	t.Parallel()
	ctx := v8.NewContext(nil)
	defer ctx.Isolate().Dispose()
	defer ctx.Close()

	math := `
	function add(a, b) {
		return a + b;
	}

	function addMore(a, b) {
		return add(a, c);
	}`

	main := `
	let a = add(3, 5);
	let b = addMore(a, 6);
	b;
	`

	ctx.RunScript(math, "math.js")
	_, err := ctx.RunScript(main, "main.js")
	if err == nil {
		t.Error("expected error but got <nil>")
		return
	}
	e, ok := err.(*v8.JSError)
	if !ok {
		t.Errorf("expected error of type JSError, got %T", err)
	}
	if e.Message != "ReferenceError: c is not defined" {
		t.Errorf("unexpected error message: %q", e.Message)
	}
	if e.Location != "math.js:7:17" {
		t.Errorf("unexpected error location: %q", e.Location)
	}
	expectedStack := `ReferenceError: c is not defined
    at addMore (math.js:7:17)
    at main.js:3:10`

	if e.StackTrace != expectedStack {
		t.Errorf("unexpected error stack trace: %q", e.StackTrace)
	}
}

func TestJSErrorFormat_forSyntaxError(t *testing.T) {
	t.Parallel()
	iso := v8.NewIsolate()
	defer iso.Dispose()
	ctx := v8.NewContext(iso)
	defer ctx.Close()

	script := `
		let x = 1;
		let y = x + ;
		let z = x + z;
	`
	_, err := ctx.RunScript(script, "xyz.js")
	jsErr := err.(*v8.JSError)
	if jsErr.StackTrace != jsErr.Message {
		t.Errorf("unexpected StackTrace %q not equal to Message %q", jsErr.StackTrace, jsErr.Message)
	}
	if jsErr.Location == "" {
		t.Errorf("missing Location")
	}

	msg := fmt.Sprintf("%+v", err)
	if msg != "SyntaxError: Unexpected token ';' (at xyz.js:3:15)" {
		t.Errorf("unexpected verbose error message: %q", msg)
	}
}
