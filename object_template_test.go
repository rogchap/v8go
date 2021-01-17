package v8go_test

import (
	"fmt"
	"testing"

	"rogchap.com/v8go"
)

func TestObjectTemplate(t *testing.T) {
	iso, _ := v8go.NewIsolate()
	global, _ := v8go.NewObjectTemplate(iso)
	val, _ := v8go.NewValue(iso, int32(23))
	global.Set("version", val, 0)
	ctx, _ := v8go.NewContext(iso, global)
	valOut, _ := ctx.RunScript("version", "test.js")
	fmt.Printf(" = %+v\n", valOut)
}
