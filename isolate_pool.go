package v8go

import (
	"context"

	"github.com/jackc/puddle"
)

// IsolatePool provides pools for isolates.
// Main advantage of using pool is the ability not to dispose
// Isolates when its not needed.
type IsolatePool struct {
	pool *puddle.Pool
}

// IsolatePoolResource wrap Isolate and add additional method
// to release resource.
type IsolatePoolResource struct {
	*Isolate
	resource *puddle.Resource
	ctx      context.Context
	cancel   context.CancelFunc
}

// Release will put resource back to pool.
func (r IsolatePoolResource) Release() {
	r.resource.Release()
	r.cancel()
}

// NewIsolatePool creates new pool of isolates.
func NewIsolatePool(poolSize int) *IsolatePool {
	constructor := func(ctx context.Context) (interface{}, error) {
		return NewIsolateContext(ctx)
	}
	destructor := func(value interface{}) {
		value.(*Isolate).Dispose()
	}
	pool := puddle.NewPool(constructor, destructor, int32(poolSize))
	return &IsolatePool{pool: pool}
}

// Acquire will get new free resource ot of pool
func (p *IsolatePool) Acquire(ctx context.Context) (*IsolatePoolResource, error) {
	res, err := p.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	iso := res.Value().(*Isolate)
	ctx, cancel := context.WithCancel(ctx)
	pr := &IsolatePoolResource{
		resource: res,
		Isolate:  iso.WithContext(ctx),
		ctx:      ctx,
		cancel:   cancel,
	}
	go func() {
		<-ctx.Done()
		iso.TerminateExecution()
	}()
	return pr, nil
}
