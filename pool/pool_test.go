package pool

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"rogchap.com/v8go"
)

func TestPool(t *testing.T) {
	p := New(3)
	var wg sync.WaitGroup

	rs := func(ctx *v8go.ExecContext) {
		defer wg.Done()
		_, err := ctx.RunScript("while(true) {}", "")
		if !strings.Contains(err.Error(), "ExecutionTerminated") {
			t.Errorf("script failed: %s", err)
		}
	}

	tp := func() {
		res, err := p.Acquire(context.Background())
		if err != nil {
			t.Error("failed to acquire pool resource")
		}
		ectx, err := v8go.NewExecContext(res.Isolate)
		if err != nil {
			t.Errorf("failed to make context: %s", err.Error())
		}
		go rs(ectx)
		time.Sleep(time.Millisecond * 100)
		res.Release()
	}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go tp()

	}

	wg.Wait()
}
