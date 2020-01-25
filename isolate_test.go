package v8go_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"rogchap.com/v8go"
)

func TestIsolateTermination(t *testing.T) {
	t.Parallel()
	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewContext(iso)
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
	v8go.NewContext(iso)
	v8go.NewContext(iso)

	hs := iso.GetHeapStatistics()

	if hs.NumberOfNativeContexts != 2 {
		t.Error("expect NumberOfNativeContexts return 2, got", hs.NumberOfNativeContexts)
	}

	if hs.NumberOfDetachedContexts != 0 {
		t.Error("expect NumberOfDetachedContexts return 0, got", hs.NumberOfDetachedContexts)
	}
}

func BenchmarkIsolateInitialization(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		vm, _ := v8go.NewIsolate()
		vm.Close() // force disposal of the VM
	}
}

func BenchmarkIsolateInitAndRun(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		vm, _ := v8go.NewIsolate()
		ctx, _ := v8go.NewContext(vm)
		ctx.RunScript(script, "main.js")
		str, _ := json.Marshal(makeObject())
		cmd := fmt.Sprintf("process(%s)", str)
		ctx.RunScript(cmd, "cmd.js")
		ctx.Close()
		vm.Close() // force disposal of the VM
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
