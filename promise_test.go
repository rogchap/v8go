// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	"rogchap.com/v8go"
)

func TestPromiseFulfilled(t *testing.T) {
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

	var thenInfo *v8go.FunctionCallbackInfo
	prom1thenVal := prom1.Then(func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		thenInfo = info
		return nil
	})
	prom1then, _ := prom1thenVal.AsPromise()
	if prom1then.State() != v8go.Pending {
		t.Errorf("unexpected state for dependent Promise, want Pending got: %v", prom1then.State())
	}
	if thenInfo != nil {
		t.Error("unexpected call of Then prior to resolving the promise")
	}

	val1, _ := v8go.NewValue(iso, "foo")
	res1.Resolve(val1)

	if s := prom1.State(); s != v8go.Fulfilled {
		t.Fatalf("unexpected state for Promise, want Fulfilled (1) got: %v", s)
	}

	if result := prom1.Result(); result.String() != val1.String() {
		t.Errorf("expected the Promise result to match the resolve value, but got: %s", result)
	}

	if thenInfo == nil {
		t.Errorf("expected Then to be called, was not")
	}
	if len(thenInfo.Args()) != 1 || thenInfo.Args()[0].String() != "foo" {
		t.Errorf("expected promise to be called with [foo] args, was: %+v", thenInfo.Args())
	}
}

func TestPromiseRejected(t *testing.T) {
	t.Parallel()

	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewContext(iso)

	res2, _ := v8go.NewPromiseResolver(ctx)
	val2, _ := v8go.NewValue(iso, "Bad Foo")
	res2.Reject(val2)

	prom2 := res2.GetPromise()
	if s := prom2.State(); s != v8go.Rejected {
		t.Fatalf("unexpected state for Promise, want Rejected (2) got: %v", s)
	}

	var thenInfo *v8go.FunctionCallbackInfo
	var then2Fulfilled, then2Rejected bool
	prom2.
		Catch(func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			thenInfo = info
			return nil
		}).
		Then2(
			func(_ *v8go.FunctionCallbackInfo) *v8go.Value {
				then2Fulfilled = true
				return nil
			},
			func(_ *v8go.FunctionCallbackInfo) *v8go.Value {
				then2Rejected = true
				return nil
			},
		)
	ctx.PerformMicrotaskCheckpoint()
	if thenInfo == nil {
		t.Fatalf("expected Then to be called on already-resolved promise, but was not")
	}
	if len(thenInfo.Args()) != 1 || thenInfo.Args()[0].String() != val2.String() {
		t.Fatalf("expected [%v], was: %+v", val2, thenInfo.Args())
	}

	if then2Fulfilled {
		t.Fatalf("unexpectedly called onFulfilled")
	}
	if !then2Rejected {
		t.Fatalf("expected call to onRejected, got none")
	}
}
