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
		{"bigbigint", bigbigint},
		{"bool", true},
		{"val", val},
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
				gbl.Set("foo", "bar", 0)
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
				foo.Set("bar", "baz", 0)
				gbl, _ := v8go.NewObjectTemplate(iso)
				gbl.Set("foo", foo, 0)
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
			ctx, _ := v8go.NewContext(iso, tt.global())
			val, err := ctx.RunScript(tt.source, "test.js")
			if err != nil {
				t.Fatalf("unexpected error runing script: %v", err)
			}
			tt.validate(t, val)
		})
	}
}
