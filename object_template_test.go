// Copyright 2020 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"math/big"
	"testing"

	"rogchap.com/v8go"
)

func TestObjectTemplate(t *testing.T) {
	t.Parallel()
	iso, _ := v8go.NewIsolate()
	defer iso.Dispose()
	obj := v8go.NewObjectTemplate(iso)

	setError := func(t *testing.T, err error) {
		if err != nil {
			t.Errorf("failed to set property: %v", err)
		}
	}

	val, _ := v8go.NewValue(iso, "bar")
	objVal := v8go.NewObjectTemplate(iso)
	bigbigint, _ := new(big.Int).SetString("36893488147419099136", 10) // larger than a single word size (64bit)
	bigbignegint, _ := new(big.Int).SetString("-36893488147419099136", 10)

	tests := [...]struct {
		name  string
		value interface{}
	}{
		{"str", "foo"},
		{"i32", int32(1)},
		{"u32", uint32(1)},
		{"i64", int64(1)},
		{"u64", uint64(1)},
		{"float64", float64(1)},
		{"bigint", big.NewInt(1)},
		{"biguint", new(big.Int).SetUint64(1 << 63)},
		{"bigbigint", bigbigint},
		{"bigbignegint", bigbignegint},
		{"bool", true},
		{"val", val},
		{"obj", objVal},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			setError(t, obj.Set(tt.name, tt.value, 0))
		})
	}
}

func TestObjectTemplate_panic_on_nil_isolate(t *testing.T) {
	t.Parallel()

	defer func() {
		if err := recover(); err == nil {
			t.Error("expected panic")
		}
	}()
	v8go.NewObjectTemplate(nil)
}

func TestGlobalObjectTemplate(t *testing.T) {
	t.Parallel()
	iso, _ := v8go.NewIsolate()
	defer iso.Dispose()
	tests := [...]struct {
		global   func() *v8go.ObjectTemplate
		source   string
		validate func(t *testing.T, val *v8go.Value)
	}{
		{
			func() *v8go.ObjectTemplate {
				gbl := v8go.NewObjectTemplate(iso)
				gbl.Set("foo", "bar")
				return gbl
			},
			"foo",
			func(t *testing.T, val *v8go.Value) {
				if !val.IsString() {
					t.Errorf("expect value %q to be of type String", val)
					return
				}
				if val.String() != "bar" {
					t.Errorf("unexpected value: %v", val)
				}
			},
		},
		{
			func() *v8go.ObjectTemplate {
				foo := v8go.NewObjectTemplate(iso)
				foo.Set("bar", "baz")
				gbl := v8go.NewObjectTemplate(iso)
				gbl.Set("foo", foo)
				return gbl
			},
			"foo.bar",
			func(t *testing.T, val *v8go.Value) {
				if val.String() != "baz" {
					t.Errorf("unexpected value: %v", val)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.source, func(t *testing.T) {
			ctx, _ := v8go.NewContext(iso, tt.global())
			val, err := ctx.RunScript(tt.source, "test.js")
			if err != nil {
				t.Fatalf("unexpected error runing script: %v", err)
			}
			tt.validate(t, val)
			ctx.Close()
		})
	}
}

func TestObjectTemplateNewInstance(t *testing.T) {
	t.Parallel()
	iso, _ := v8go.NewIsolate()
	defer iso.Dispose()
	tmpl := v8go.NewObjectTemplate(iso)
	if _, err := tmpl.NewInstance(nil); err == nil {
		t.Error("expected error but got <nil>")
	}

	tmpl.Set("foo", "bar")
	ctx, _ := v8go.NewContext(iso)
	defer ctx.Close()
	obj, _ := tmpl.NewInstance(ctx)
	if foo, _ := obj.Get("foo"); foo.String() != "bar" {
		t.Errorf("unexpected value for object property: %v", foo)
	}

}
