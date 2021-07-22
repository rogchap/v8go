// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"fmt"
	"testing"

	"github.com/Shopify/v8go"
)

func TestJSONParse(t *testing.T) {
	t.Parallel()

	if _, err := v8go.JSONParse(nil, "{}"); err == nil {
		t.Error("expected error but got <nil>")
	}
	ctx, _ := v8go.NewContext()
	_, err := v8go.JSONParse(ctx, "{")
	if err == nil {
		t.Error("expected error but got <nil>")
		return
	}

	if _, ok := err.(*v8go.JSError); !ok {
		t.Errorf("expected error to be of type JSError, got: %T", err)
	}
}

func TestJSONStringify(t *testing.T) {
	t.Parallel()

	ctx, _ := v8go.NewContext()
	if _, err := v8go.JSONStringify(ctx, nil); err == nil {
		t.Error("expected error but got <nil>")
	}
}

func ExampleJSONParse() {
	ctx, _ := v8go.NewContext()
	val, _ := v8go.JSONParse(ctx, `{"foo": "bar"}`)
	fmt.Println(val)
	// Output:
	// [object Object]
}

func ExampleJSONStringify() {
	ctx, _ := v8go.NewContext()
	val, _ := v8go.JSONParse(ctx, `{
		"a": 1,
		"b": "foo"
	}`)
	jsonStr, _ := v8go.JSONStringify(ctx, val)
	fmt.Println(jsonStr)
	// Output:
	// {"a":1,"b":"foo"}
}
