package v8go_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"rogchap.com/v8go"
)

func TestFunctionTemplate(t *testing.T) {
	//TODO: write proper tests

	iso, _ := v8go.NewIsolate()
	global, _ := v8go.NewObjectTemplate(iso)
	printfn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		fmt.Printf("Called print(): %+v\n", args)
		val, _ := v8go.NewValue(iso, int32(len(args)))
		return val
	})
	logfn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) (rtn *v8go.Value) {
		fmt.Printf("Called log(): %+v\n", info.Args())
		return
	})
	th := &thing{}
	headfn, _ := v8go.NewFunctionTemplate(iso, th.cb)

	global.Set("print", printfn)
	global.Set("log", logfn)
	global.Set("head", headfn)
	ctx, _ := v8go.NewContext(iso, global)
	val, _ := ctx.RunScript("head();log(print);print('stuff', 'more', 4, 2n)", "")
	fmt.Printf("val = %+v\n", val)
}

type thing struct {
	c context.Context
	d *http.Client
}

func (t *thing) cb(info *v8go.FunctionCallbackInfo) *v8go.Value {
	t.d = http.DefaultClient
	go t.d.Head("http://rogchap.com")
	return nil
}
