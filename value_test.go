package v8go_test

import (
	"reflect"
	"runtime"
	"testing"

	"rogchap.com/v8go"
)

func TestValueString(t *testing.T) {
	t.Parallel()
	ctx, _ := v8go.NewContext(nil)
	var tests = [...]struct {
		name   string
		source string
		out    string
	}{
		{"Number", `13 * 2`, "26"},
		{"String", `"string"`, "string"},
		{"Object", `let obj = {}; obj`, "[object Object]"},
		{"Function", `let fn = function(){}; fn`, "function(){}"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, _ := ctx.RunScript(tt.source, "test.js")
			str := result.String()
			if str != tt.out {
				t.Errorf("unespected result: expected %q, got %q", tt.out, str)
			}
		})
	}
}

func TestValueIsXXX(t *testing.T) {
	t.Parallel()
	iso, _ := v8go.NewIsolate()
	var tests = []struct {
		source string
		assert func(*v8go.Value) bool
	}{
		{"", (*v8go.Value).IsUndefined},
		{"let v; v", (*v8go.Value).IsUndefined},
		{"let v = null; v", (*v8go.Value).IsNull},
		{"let v; v", (*v8go.Value).IsNullOrUndefined},
		{"let v = null; v", (*v8go.Value).IsNullOrUndefined},
		{"let v = true; v", (*v8go.Value).IsTrue},
		{"let v = false; v", (*v8go.Value).IsFalse},
		{`"double quote"`, (*v8go.Value).IsString},
		{"'single quote'", (*v8go.Value).IsString},
		{"`string litral`", (*v8go.Value).IsString},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.source, func(t *testing.T) {
			t.Parallel()
			ctx, _ := v8go.NewContext(iso)
			val, err := ctx.RunScript(tt.source, "test.js")
			if err != nil {
				t.Fatalf("failed to run script: %v", err)
			}
			if !tt.assert(val) {
				t.Errorf("value is false for %s", runtime.FuncForPC(reflect.ValueOf(tt.assert).Pointer()).Name())
			}
		})
	}
}
