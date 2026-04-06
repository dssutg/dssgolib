// Package brokers implements a broker pattern for
// subscription and message broadcasting.
package brokers

// Broker represents a broadcast system that tracks
// subscribers and sends them updates.
type Broker[T any] struct {
	stopCh    chan struct{}
	publishCh chan T
	attachCh  chan chan T
	detachCh  chan chan T
	unsubCh   chan chan T
}

// New creates and returns a new Broker instance with the
// provided channel capacity.
func New[T any](capacity int) *Broker[T] {
	return &Broker[T]{
		stopCh:    make(chan struct{}),
		publishCh: make(chan T, capacity),
		attachCh:  make(chan chan T, capacity),
		detachCh:  make(chan chan T, capacity),
		unsubCh:   make(chan chan T, capacity),
	}
}

// Start handles registering, unregistering, and
// broadcasting messages to subscribers.
func (b *Broker[T]) Start() {
	subs := map[chan T]struct{}{}
	for {
		select {
		case <-b.stopCh:
			for sub := range subs {
				close(sub)
			}
			clear(subs)
			return
		case sub := <-b.attachCh:
			subs[sub] = struct{}{}
		case sub := <-b.detachCh:
			delete(subs, sub)
		case sub := <-b.unsubCh:
			if _, ok := subs[sub]; ok {
				delete(subs, sub)
				close(sub)
			}
		case msg := <-b.publishCh:
			for sub := range subs {
				// sub is buffered, use non-blocking send
				// to protect the broker:
				select {
				case sub <- msg:
				default:
				}
			}
		}
	}
}

// Stop terminates the broker. It closes all subscriber channels and
// clears the subscriber queue.
func (b *Broker[T]) Stop() {
	close(b.stopCh)
}

// Attach adds the subscriber to the subscriber queue.
func (b *Broker[T]) Attach(sub chan T) {
	b.attachCh <- sub
}

// Detach removes the subscriber from the subscriber queue.
// If you want to also close the channel, you must use
// the Unsubscribe method instead. if there's not such
// a subscriber, this function is no-op.
func (b *Broker[T]) Detach(sub chan T) {
	b.detachCh <- sub
}

// Subscribe creates and returns a new subscriber's channel with the given
// capacity, then adds it the subscriber queue.
func (b *Broker[T]) Subscribe(capacity int) chan T {
	sub := make(chan T, capacity)
	b.Attach(sub)
	return sub
}

// Unsubscribe removes the subscriber's channel from the
// broker queue, then closes the channel.
// If the channel is nil, this function is no-op.
// if there's not such a subscriber, this function is no-op.
func (b *Broker[T]) Unsubscribe(sub chan T) {
	if sub == nil {
		return
	}
	// Don't close sub here to avoid a race where the broker
	// could send to a closed channel.
	b.unsubCh <- sub
}

// Publish sends a message to all subscribers.
func (b *Broker[T]) Publish(msg T) {
	b.publishCh <- msg
}
