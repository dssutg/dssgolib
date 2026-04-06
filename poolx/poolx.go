package poolx

import "sync"

// Pool is the type-safe version of sync.pool.
type Pool[T any] struct {
	pool sync.Pool
}

// Make creates a new pool with the provided maker for each new object.
func Make[T any](maker func() *T) Pool[T] {
	return Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return maker()
			},
		},
	}
}

// New returns an object from the pool,
func (p *Pool[T]) New() *T {
	return p.pool.Get().(*T)
}

// Free zeroes the object, then gets it back to the pool.
func (p *Pool[T]) Free(v *T) {
	var zero T
	*v = zero
	p.pool.Put(v)
}
