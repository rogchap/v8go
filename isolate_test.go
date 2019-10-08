package v8go_test

import (
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
