package v8go_test

import (
	"fmt"
	"testing"

	"rogchap.com/v8go"
)

func TestFunctionTemplate(t *testing.T) {
	//TODO: write proper tests

	iso, _ := v8go.NewIsolate()
	global, _ := v8go.NewObjectTemplate(iso)
	printfn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		fmt.Printf("Called print(): %+v\n", info.Args())
		return nil
	})
	logfn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		fmt.Printf("Called log(): %+v\n", info.Args())
		return nil
	})

	global.Set("print", printfn)
	global.Set("log", logfn)
	ctx, _ := v8go.NewContext(iso, global)
	ctx.RunScript("log();print('stuff', 'more', 4, 2n)", "")

}
