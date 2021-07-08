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
	_, err := v8go.NewObjectTemplate(nil)
	if err == nil {
		t.Fatal("expected error but got <nil>")
	}
	iso, _ := v8go.NewIsolate()
	obj, err := v8go.NewObjectTemplate(iso)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	setError := func(t *testing.T, err error) {
		if err != nil {
			t.Errorf("failed to set property: %v", err)
		}
	}

	val, _ := v8go.NewValue(iso, "bar")
	objVal, _ := v8go.NewObjectTemplate(iso)
	bigbigint, _ := new(big.Int).SetString("36893488147419099136", 10) // larger than a single word size (64bit)
	bigbignegint, _ := new(big.Int).SetString("-36893488147419099136", 10)

	valtemp := func(x interface{}) v8go.Templater {
		v, err := v8go.NewValueTemplate(iso, x)
		if err != nil {
			t.Errorf("failed to make property: %v", err)
		}
		return v
	}

	tests := [...]struct {
		name  string
		value v8go.Templater
	}{
		{"str", valtemp("foo")},
		{"i32", valtemp(int32(1))},
		{"u32", valtemp(uint32(1))},
		{"i64", valtemp(int64(1))},
		{"u64", valtemp(uint64(1))},
		{"float64", valtemp(float64(1))},
		{"bigint", valtemp(big.NewInt(1))},
		{"biguint", valtemp(new(big.Int).SetUint64(1 << 63))},
		{"bigbigint", valtemp(bigbigint)},
		{"bigbignegint", valtemp(bigbignegint)},
		{"bool", valtemp(true)},
		{"val", valtemp(val)},
		{"obj", objVal},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			setError(t, obj.Set(tt.name, tt.value, 0))
		})
	}
}

func TestGlobalObjectTemplate(t *testing.T) {
	t.Parallel()
	iso, _ := v8go.NewIsolate()
	tests := [...]struct {
		global   func() *v8go.ObjectTemplate
		source   string
		validate func(t *testing.T, val *v8go.Value)
	}{
		{
			func() *v8go.ObjectTemplate {
				gbl, _ := v8go.NewObjectTemplate(iso)
				v, _ := v8go.NewValueTemplate(iso, "bar")
				gbl.Set("foo", v)
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
				foo, _ := v8go.NewObjectTemplate(iso)
				v, _ := v8go.NewValueTemplate(iso, "baz")
				foo.Set("bar", v)
				gbl, _ := v8go.NewObjectTemplate(iso)
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
			t.Parallel()
			ctx, _ := v8go.NewExecContext(iso, tt.global())
			val, err := ctx.RunScript(tt.source, "test.js")
			if err != nil {
				t.Fatalf("unexpected error runing script: %v", err)
			}
			tt.validate(t, val)
		})
	}
}

func TestObjectTemplateNewInstance(t *testing.T) {
	t.Parallel()
	iso, _ := v8go.NewIsolate()
	tmpl, _ := v8go.NewObjectTemplate(iso)
	if _, err := tmpl.GetObject(nil); err == nil {
		t.Error("expected error but got <nil>")
	}

	v, _ := v8go.NewValueTemplate(iso, "bar")
	tmpl.Set("foo", v)
	ctx, _ := v8go.NewExecContext(iso)
	obj, _ := tmpl.GetObject(ctx)
	if foo, _ := obj.Get("foo"); foo.String() != "bar" {
		t.Errorf("unexpected value for object property: %v", foo)
	}

}
