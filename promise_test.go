// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	"github.com/Shopify/v8go"
)

func TestPromise(t *testing.T) {
	t.Parallel()

	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewContext(iso)
	if _, err := v8go.NewPromiseResolver(nil); err == nil {
		t.Error("expected error with <nil> Context")
	}

	res1, _ := v8go.NewPromiseResolver(ctx)
	prom1 := res1.GetPromise()
	if s := prom1.State(); s != v8go.Pending {
		t.Errorf("unexpected state for Promise, want Pending (0) got: %v", s)
	}

	val1, _ := v8go.NewValue(iso, "foo")
	res1.Resolve(val1)

	if s := prom1.State(); s != v8go.Fulfilled {
		t.Fatalf("unexpected state for Promise, want Fulfilled (1) got: %v", s)
	}

	if result := prom1.Result(); result.String() != val1.String() {
		t.Errorf("expected the Promise result to match the resolve value, but got: %s", result)
	}

	res2, _ := v8go.NewPromiseResolver(ctx)
	val2, _ := v8go.NewValue(iso, "Bad Foo")
	res2.Reject(val2)

	prom2 := res2.GetPromise()
	if s := prom2.State(); s != v8go.Rejected {
		t.Fatalf("unexpected state for Promise, want Rejected (2) got: %v", s)
	}
}
