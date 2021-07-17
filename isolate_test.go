// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"rogchap.com/v8go"
)

func TestIsolateTermination(t *testing.T) {
	t.Parallel()
	iso, _ := v8go.NewIsolate()
	iso = iso.WithContext(context.Background())
	ctx, _ := v8go.NewExecContext(iso)
	//	ctx2, _ := v8go.NewContext(iso)

	err := make(chan error, 1)

	go func() {
		_, e := ctx.RunScript(`while (true) { }`, "forever.js")
		err <- e
	}()

	go func() {
		// [RC] find a better way to know when a script has started execution
		time.Sleep(time.Millisecond)
		iso.TerminateExecution()
	}()

	if e := <-err; e == nil || !strings.HasPrefix(e.Error(), "ExecutionTerminated") {
		t.Errorf("unexpected error: %v", e)
	}
}

func TestGetHeapStatistics(t *testing.T) {
	t.Parallel()
	iso, _ := v8go.NewIsolate()
	v8go.NewExecContext(iso)
	v8go.NewExecContext(iso)

	hs := iso.GetHeapStatistics()

	if hs.NumberOfNativeContexts != 3 {
		t.Error("expect NumberOfNativeContexts return 3, got", hs.NumberOfNativeContexts)
	}

	if hs.NumberOfDetachedContexts != 0 {
		t.Error("expect NumberOfDetachedContexts return 0, got", hs.NumberOfDetachedContexts)
	}
}

func TestCallbackRegistry(t *testing.T) {
	t.Parallel()

	iso, _ := v8go.NewIsolate()
	iso = iso.WithContext(context.Background())
	cb := func(*v8go.FunctionCallbackInfo) (v8go.Valuer, error) { return nil, nil }

	cb0 := iso.GetCallback(0)
	if cb0 != nil {
		t.Error("expected callback function to be <nil>")
	}
	ref1 := iso.RegisterCallback(cb)
	if ref1 != 1 {
		t.Errorf("expected callback ref == 1, got %d", ref1)
	}
	cb1 := iso.GetCallback(1)
	if fmt.Sprintf("%p", cb1) != fmt.Sprintf("%p", cb) {
		t.Errorf("unexpected callback function; want %p, got %p", cb, cb1)
	}
}

func TestIsolateDispose(t *testing.T) {
	t.Parallel()

	iso, _ := v8go.NewIsolate()
	if iso.GetHeapStatistics().TotalHeapSize == 0 {
		t.Error("Isolate incorrectly allocated")
	}

	iso.Dispose()
	// noop when called multiple times
	iso.Dispose()

	if iso.GetHeapStatistics().TotalHeapSize != 0 {
		t.Error("Isolate not disposed correctly")
	}
}

func TestIsolateGarbageCollection(t *testing.T) {
	t.Parallel()

	iso, _ := v8go.NewIsolate()
	val, _ := v8go.NewValue(iso, "some string")
	fmt.Println(val.String())

	tmpl, _ := v8go.NewObjectTemplate(iso)
	v, _ := v8go.NewValueTemplate(iso, "bar")

	tmpl.Set("foo", v)
	v8go.NewExecContext(iso, tmpl)

	iso.Dispose()

	runtime.GC()

	time.Sleep(time.Second)
}

func BenchmarkIsolateInitialization(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		vm, _ := v8go.NewIsolate()
		vm.Dispose() // force disposal of the VM
	}
}

func BenchmarkIsolateInitAndRun(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		vm, _ := v8go.NewIsolate()
		ctx, _ := v8go.NewExecContext(vm)
		ctx.RunScript(script, "main.js")
		str, _ := json.Marshal(makeObject())
		cmd := fmt.Sprintf("process(%s)", str)
		ctx.RunScript(cmd, "cmd.js")
		ctx.Close()
		vm.Dispose() // force disposal of the VM
	}
}

const script = `
	const process = (record) => {
		const res = [];
		for (let [k, v] of Object.entries(record)) {
			res.push({
				name: k,
				value: v,
			});
		}
		return JSON.stringify(res);
	};
`

func makeObject() interface{} {
	return map[string]interface{}{
		"a": rand.Intn(1000000),
		"b": "AAAABBBBAAAABBBBAAAABBBBAAAABBBBAAAABBBB",
	}
}

func TestDisposal(t *testing.T) {
	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewExecContext(iso)
	m := make(chan struct{})
	go func() {
		m <- struct{}{}
		ctx.RunScript("while(true) {}", "")
	}()
	<-m
	err := iso.Dispose()
	if err != v8go.ErrIsolateInUse {
		t.Error("error should be IsolateInUse")
	}
	iso.TerminateExecutionWithLock()
	err = iso.Dispose()
	if err != nil {
		t.Error("not possible to dispose")
	}
}

func TestExecReturnValues(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	iso, _ := v8go.NewIsolateContext(ctx)

	var wg1 sync.WaitGroup
	wg1.Add(10)

	var vv []string

	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg1.Done()
			err := iso.Exec(func(iso *v8go.Isolate) error {
				v, _ := v8go.NewExecContext(iso)
				vx, err := v.RunScript(fmt.Sprintf("%d", i), "")
				vv = append(vv, vx.String())
				return err
			})
			if err != nil {
				t.Errorf("error from safe: %s", err.Error())
			}
		}(i)
	}
	// All go routines has been spawned
	wg1.Wait()
	sort.Strings(vv)
	joined := strings.Join(vv, ",")
	if joined != "0,1,2,3,4,5,6,7,8,9" {
		t.Errorf("invalid result returned: %s expects %s", joined, "0,1,2,3,4,5,6,7,8,9")
	}
}
