package v8go_test

import (
	"fmt"

	"rogchap.com/v8go"
)

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
