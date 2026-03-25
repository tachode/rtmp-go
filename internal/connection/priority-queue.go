package connection

import "sync"

type PriorityQueue[T any] struct {
	done     chan struct{}
	once     sync.Once
	ready    chan bool
	priority []chan T
}

// New creates a new Queue.
func NewPriorityQueue[T any](priorityCount int, queueLength int) *PriorityQueue[T] {
	q := &PriorityQueue[T]{
		done:     make(chan struct{}),
		ready:    make(chan bool, priorityCount*queueLength),
		priority: make([]chan T, priorityCount),
	}
	for i := range q.priority {
		q.priority[i] = make(chan T, queueLength)
	}
	return q
}

// Enqueue adds a value to the queue. Returns an error if the queue is closed.
func (q *PriorityQueue[T]) Enqueue(v T, priority int) error {
	if priority < 0 || priority >= len(q.priority) {
		priority = len(q.priority) - 1
	}
	select {
	case q.priority[priority] <- v:
	case <-q.done:
		return ErrQueueClosed
	}
	select {
	case q.ready <- true:
	case <-q.done:
		return ErrQueueClosed
	}
	return nil
}

// Dequeue removes a value from the queue. Returns the zero value and false
// if the queue has been closed.
func (q *PriorityQueue[T]) Dequeue() (T, bool) {
	select {
	case <-q.ready:
	case <-q.done:
		return *new(T), false
	}
	for _, in := range q.priority {
		select {
		case v := <-in:
			return v, true
		default:
		}
	}
	return *new(T), false
}

func (q *PriorityQueue[T]) Length() int {
	return len(q.ready)
}

// Close signals the queue to stop. Safe to call multiple times.
func (q *PriorityQueue[T]) Close() {
	q.once.Do(func() {
		close(q.done)
	})
}
