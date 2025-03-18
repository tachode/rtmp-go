package connection_test

import (
	"testing"

	"github.com/tachode/rtmp-go/internal/connection"
)

func TestQueue_EnqueueDequeue(t *testing.T) {
	q := connection.NewPriorityQueue[int](3, 10)

	// Enqueue items with different priorities
	q.Enqueue(1, 0)
	q.Enqueue(2, 1)
	q.Enqueue(3, 2)

	// Dequeue items and check order
	if v := q.Dequeue(); v != 1 {
		t.Errorf("expected 1, got %d", v)
	}
	if v := q.Dequeue(); v != 2 {
		t.Errorf("expected 2, got %d", v)
	}
	if v := q.Dequeue(); v != 3 {
		t.Errorf("expected 3, got %d", v)
	}
}

func TestQueue_PriorityOrder(t *testing.T) {
	q := connection.NewPriorityQueue[int](3, 10)

	// Enqueue items with different priorities
	q.Enqueue(1, 2) // Lowest priority
	q.Enqueue(2, 0) // Highest priority
	q.Enqueue(3, 1) // Medium priority

	// Dequeue items and check priority order
	if v := q.Dequeue(); v != 2 {
		t.Errorf("expected 2, got %d", v)
	}
	if v := q.Dequeue(); v != 3 {
		t.Errorf("expected 3, got %d", v)
	}
	if v := q.Dequeue(); v != 1 {
		t.Errorf("expected 1, got %d", v)
	}
}

func TestQueue_InvalidPriority(t *testing.T) {
	q := connection.NewPriorityQueue[int](3, 10)

	// Enqueue with invalid priority (negative)
	q.Enqueue(1, -1)

	// Enqueue with invalid priority (too high)
	q.Enqueue(2, 10)

	// Dequeue items and check they are enqueued in the lowest priority
	if v := q.Dequeue(); v != 1 {
		t.Errorf("expected 1, got %d", v)
	}
	if v := q.Dequeue(); v != 2 {
		t.Errorf("expected 2, got %d", v)
	}
}
