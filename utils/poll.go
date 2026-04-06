package utils

import (
	"context"
	"time"
)

// Poll immediately calls fn, then waits for the interval. When ctx is done,
// the goroutine returns.
func Poll(ctx context.Context, interval time.Duration, fn func()) {
	go func() {
		fn() // immediate
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fn()
			}
		}
	}()
}
