package pool

import (
	"context"
	"runtime"

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
}

// Release will put resource back to pool and
// terminates the execution.
func (r Resource) Release() {
	r.Isolate.TerminateExecution()
	r.resource.Release()
}

// New creates new pool of isolates.
func New(poolSize int) *Pool {
	constructor := func(ctx context.Context) (interface{}, error) {
		return v8go.NewIsolateContext(ctx)
	}
	destructor := func(value interface{}) {
		iso := value.(*v8go.Isolate)
		iso.TerminateExecutionWithLock()
		err := iso.Dispose()
		if err != nil {
			panic(err)
		}
	}
	pool := puddle.NewPool(constructor, destructor, int32(poolSize))
	return &Pool{pool: pool}
}

// Default is the default pool based on number of cpus.
var Default = New(runtime.NumCPU())

// Acquire will get new free resource ot of pool.
func (p *Pool) Acquire(ctx context.Context) (*Resource, error) {
	res, err := p.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	iso := res.Value().(*v8go.Isolate)
	pr := &Resource{
		resource: res,
		Isolate:  iso.WithContext(ctx),
	}
	return pr, nil
}

// Close will close the pool and dipose all vms.
func (p *Pool) Close() {
	p.pool.Close()
}
