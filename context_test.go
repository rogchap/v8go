package v8go_test

import (
	"testing"

	"rogchap.com/v8go"
)

func TestRunScriptStringer(t *testing.T) {
	t.Parallel()
	ctx, _ := v8go.NewContext(nil)
	var tests = [...]struct {
		name   string
		source string
		out    string
	}{
		{"Addition", `2 + 2`, "4"},
		{"Multiplication", `13 * 2`, "26"},
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
