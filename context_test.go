package v8go_test

import (
	"testing"

	"rogchap.com/v8go"
)

func TestContextExec(t *testing.T) {
	t.Parallel()
	ctx, _ := v8go.NewContext(nil)
	ctx.RunScript(`function add(a, b) { return a + b }`, "add.js")
	val, _ := ctx.RunScript(`add(3, 4)`, "main.js")
	rtn := val.String()
	if rtn != "7" {
		t.Errorf("script returned an unexpected value: expected %q, got %q", "7", rtn)
	}

	_, err := ctx.RunScript(`add`, "func.js")
	if err != nil {
		t.Errorf("function should be defined: %v", err)
	}

	/*iso, _ := ctx.Isolate()
	ctx2, _ := v8go.NewContext(iso)
	val, _ = ctx2.RunScript(`add`, "ctx2.js")
	rtn = val.String()
	if rtn != "undefined" {
		t.Errorf("%q", rtn)
	}*/
}
