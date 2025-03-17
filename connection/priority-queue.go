package connection

type PriorityQueue[T any] struct {
	ready    chan bool
	priority []chan T
}

// New creates a new Queue.
func NewPriorityQueue[T any](priorityCount int, queueLength int) *PriorityQueue[T] {
	q := &PriorityQueue[T]{
		ready:    make(chan bool, priorityCount*queueLength),
		priority: make([]chan T, priorityCount),
	}
	for i := range q.priority {
		q.priority[i] = make(chan T, queueLength)
	}
	return q
}

// Enqueue adds a value to the queue.
func (q *PriorityQueue[T]) Enqueue(v T, priority int) {
	if priority < 0 || priority >= len(q.priority) {
		priority = len(q.priority) - 1
	}
	q.priority[priority] <- v
	q.ready <- true
}

// Dequeue removes a value from the queue.
func (q *PriorityQueue[T]) Dequeue() T {
	open := <-q.ready
	if !open {
		return *new(T)
	}
	for _, in := range q.priority {
		select {
		case v := <-in:
			return v
		default:
		}
	}
	return *new(T)
}

func (q *PriorityQueue[T]) Length() int {
	return len(q.ready)
}

func (q *PriorityQueue[T]) Close() {
	close(q.ready)
	for _, priority := range q.priority {
		close(priority)
	}
}
