package workermanager

import (
	"context"
	"sync"
)

type WorkerManager[T comparable] struct {
	mu     sync.Mutex
	cancel map[T]context.CancelFunc
	wg     sync.WaitGroup
}

func NewWorkerManager[T comparable]() *WorkerManager[T] {
	return &WorkerManager[T]{cancel: make(map[T]context.CancelFunc)}
}

func (m *WorkerManager[T]) Start(id T, work func(ctx context.Context, c context.CancelFunc)) {
	m.mu.Lock()
	if _, ok := m.cancel[id]; ok {
		m.mu.Unlock()
		return // already started
	}
	ctx, c := context.WithCancel(context.Background())
	m.cancel[id] = c
	m.mu.Unlock()

	m.wg.Go(func() {
		defer func() {
			// remove from map when exited
			m.mu.Lock()
			delete(m.cancel, id)
			m.mu.Unlock()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				work(ctx, c)
			}
		}
	})
}

func (m *WorkerManager[T]) Stop(id T) {
	m.mu.Lock()
	c, ok := m.cancel[id]
	m.mu.Unlock()
	if ok {
		c() // signal that worker to stop
	}
}

func (m *WorkerManager[T]) StopMany(ids []T) {
	m.mu.Lock()
	for _, id := range ids {
		if c, ok := m.cancel[id]; ok {
			c()
		}
	}
	m.mu.Unlock()
}

func (m *WorkerManager[T]) StopAll() {
	m.mu.Lock()
	for _, c := range m.cancel {
		c()
	}
	m.mu.Unlock()
}

func (m *WorkerManager[T]) Wait() {
	m.wg.Wait()
}
