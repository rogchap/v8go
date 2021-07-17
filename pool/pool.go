package pool

import (
	"context"

	"github.com/jackc/puddle"
	"rogchap.com/v8go"
)

// Pool provides pools for isolates.
// Main advantage of using pool is the ability not to dispose
// Isolates when its not needed.
type Pool struct {
	pool *puddle.Pool
}

// Resource wrap Isolate and add additional method
// to release resource.
type Resource struct {
	*v8go.Isolate
	resource *puddle.Resource
	ctx      context.Context
	cancel   context.CancelFunc
}

// Release will put resource back to pool.
func (r Resource) Release() {
	r.resource.Release()
	r.cancel()
}

// New creates new pool of isolates.
func New(poolSize int) *Pool {
	constructor := func(ctx context.Context) (interface{}, error) {
		return v8go.NewIsolateContext(ctx)
	}
	destructor := func(value interface{}) {
		value.(*v8go.Isolate).Dispose()
	}
	pool := puddle.NewPool(constructor, destructor, int32(poolSize))
	return &Pool{pool: pool}
}

// Acquire will get new free resource ot of pool
func (p *Pool) Acquire(ctx context.Context) (*Resource, error) {
	res, err := p.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	iso := res.Value().(*v8go.Isolate)
	ctx, cancel := context.WithCancel(ctx)
	pr := &Resource{
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
