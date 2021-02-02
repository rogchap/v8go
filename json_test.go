package v8go_test

import (
	"fmt"

	"rogchap.com/v8go"
)

func ExampleJSON_parse() {
	ctx, _ := v8go.NewContext()
	val, _ := v8go.JSONParse(ctx, `{"foo": "bar"}`)
	fmt.Println(val)
	// Output:
	// [object Object]
}

func ExampleJSON_stringify() {
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
