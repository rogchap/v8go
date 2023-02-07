package v8go_test

import (
	"encoding/json"
	"fmt"
	"testing"

	v8 "rogchap.com/v8go"
)

func TestHeapSnapshot(t *testing.T) {
	t.Parallel()
	iso := v8.NewIsolate()
	heapProfiler := v8.NewHeapProfiler(iso)
	defer iso.Dispose()
	printfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		fmt.Printf("%v", info.Args())
		return nil
	})
	global := v8.NewObjectTemplate(iso)
	global.Set("print", printfn)
	ctx := v8.NewContext(iso, global)
	ctx.RunScript("print('foo')", "print.js")

	str, err := heapProfiler.TakeHeapSnapshot()
	if err != nil {
		t.Errorf("expected nil but got error: %v", err)
	}

	var snapshot map[string]interface{}
	err = json.Unmarshal([]byte(str), &snapshot)
	if err != nil {
		t.Fatal(err)
	}
}
