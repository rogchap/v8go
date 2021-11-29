// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	v8 "rogchap.com/v8go"
)

func TestIntlSupport(t *testing.T) {
	t.Parallel()

	ctx := v8.NewContext(nil)
	iso := ctx.Isolate()
	defer iso.Dispose()
	defer ctx.Close()

	v, err := ctx.RunScript("typeof Intl === 'object'", "test.js")
	fatalIf(t, err)
	if !v.Boolean() {
		t.Fatalf("expected value to be true, but was false")
	}

	v, err = ctx.RunScript("new Intl.NumberFormat()", "test.js")
	fatalIf(t, err)
	if v.String() != "[object Intl.NumberFormat]" {
		t.Fatalf("expected value to be [object Intl.NumberFormat], but was %v", v)
	}

	v, err = ctx.RunScript("new Intl.DateTimeFormat('es', { month: 'long' }).format(new Date(9E8))", "test.js")
	fatalIf(t, err)
	if v.String() != "enero" {
		t.Fatalf("expected value to be enero, but was %v", v)
	}

	// Example from the node.js documentation:
	// https://github.com/nodejs/node/blob/2e2a6fecd9b1aaffcb932fcc415439f359c84fdd/doc/api/intl.md?plain=1#L175-L183
	script := `const hasFullICU = (() => {
try {
	const january = new Date(9e8);
	const spanish = new Intl.DateTimeFormat('es', { month: 'long' });
	return spanish.format(january) === 'enero';
} catch (err) {
	return false;
}
})(); hasFullICU`

	v, err = ctx.RunScript(script, "test.js")
	fatalIf(t, err)
	if !v.Boolean() {
		t.Fatalf("expected value to be true, but was %v", v)
	}

	v, err = ctx.RunScript("var number = 123456.789; new Intl.NumberFormat('de-DE').format(number)", "test.js")
	fatalIf(t, err)
	if v.String() != "123.456,789" {
		t.Fatalf("expected value to be %v, but was %v", "123.456,789", v)
	}
}
