package v8go_test

import (
	"fmt"
	"testing"

	"rogchap.com/v8go"
)

func TestObjectSet(t *testing.T) {

}

func TestObjectGet(t *testing.T) {

}

func TestObjectHas(t *testing.T) {

}

func TestObjectDelete(t *testing.T) {

}

func ExampleObject_global() {
	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewContext(iso)
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
