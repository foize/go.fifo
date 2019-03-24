package fifo

import (
	"sync"
)

// fifo queue
type Queue struct {
	l     sync.Mutex
	queue *UnsafeQueue
}

// NewUnsafeQueue creates a new and empty *fifo.UnsafeQueue
func NewQueue() (q *Queue) {
	q = &Queue{
		l:     sync.Mutex{},
		queue: NewUnsafeQueue(),
	}
	return q
}

// Return the number of items in the queue
func (q *Queue) Len() (length int) {
	// locking to make UnsafeQueue thread-safe
	q.l.Lock()
	c := q.queue.Len()
	q.l.Unlock()

	return c
}

// Add an item to the end of the queue
func (q *Queue) Add(item interface{}) {
	// locking to make UnsafeQueue thread-safe
	q.l.Lock()
	q.queue.Add(item)
	q.l.Unlock()
}

func (q *Queue) AddList(items []interface{}) {
	q.l.Lock()
	q.queue.AddList(items)
	q.l.Unlock()
}

// Remove the item at the head of the queue and return it.
// Returns nil when there are no items left in queue.
func (q *Queue) Next() (item interface{}) {
	// locking to make UnsafeQueue thread-safe
	q.l.Lock()
	i := q.queue.Next()
	q.l.Unlock()

	return i
}

// Returns the next N elements from the queue
// In case of not enough elements, returns the elements that are available
func (q *Queue) NextN(n int) []interface{} {
	q.l.Lock()
	i := q.queue.NextN(n)
	q.l.Unlock()

	return i
}

// Reads the item at the head of the queue without removing it
// Returns nil when there are no items left in queue
func (q *Queue) Peek() (item interface{}) {
	q.l.Lock()
	i := q.queue.Peek()
	q.l.Unlock()

	return i
}
