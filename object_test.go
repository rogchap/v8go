package v8go_test

import (
	"fmt"

	"rogchap.com/v8go"
)

func ExampleObject_global() {
	iso, _ := v8go.NewIsolate()
	obj, _ := v8go.NewObjectTemplate(iso)
	obj.Set("version", "v1.0.0")
	ctx, _ := v8go.NewContext(iso, obj)
	global := ctx.Global()

	console, _ := v8go.NewObjectTemplate(iso)
	logfn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		fmt.Println(info.Args()[0])
		return nil
	})
	console.Set("log", logfn)
	consoleObj, _ := console.NewInstance(ctx)

	global.Set("console", consoleObj)
	ctx.RunScript("console.log('foo')", "")
	// Output:
	// foo
}
