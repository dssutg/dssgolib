package debounce

import (
	"sync"
	"time"
)

// NewUpdateDebouncer immediately calls f() once when invoked. If subsequent
// calls occur during the wait period, schedules one extra call of f() at the
// end of the waiting period.
func NewUpdateDebouncer(wait time.Duration, f func()) func() {
	var mu sync.Mutex
	var t *time.Timer
	var extra bool

	return func() {
		mu.Lock()
		defer mu.Unlock()

		// If timer exists, we are in the wait period:
		if t != nil {
			extra = true
			return
		}

		// No timer running; process the update immediately:
		extra = false
		// Call f() outside of lock to avoid long lock holding if f() is slow.
		mu.Unlock()
		f()
		mu.Lock()

		// Create a timer to check for extra update after the wait period.
		t = time.AfterFunc(wait, func() {
			var callExtra bool

			mu.Lock()
			if extra {
				// Capture that we need to call f()
				callExtra = true
				// Clear the extra flag as we've "satisfied" the extra update.
				extra = false
			}
			// Clear timer before releasing the lock.
			t = nil
			mu.Unlock()

			// Now call f() outside the lock, if needed.
			if callExtra {
				f()
			}
		})
	}
}
