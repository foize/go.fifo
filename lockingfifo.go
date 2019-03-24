package fifo

import (
	"sync"
)

// fifo queue
type LockingQueue struct {
	l     sync.Mutex
	queue *Queue
}

// NewQueue creates a new and empty *fifo.Queue
func NewLockingQueue() (q *LockingQueue) {
	q = &LockingQueue{
		l:     sync.Mutex{},
		queue: NewQueue(),
	}
	return q
}

// Return the number of items in the queue
func (q *LockingQueue) Len() (length int) {
	// locking to make Queue thread-safe
	q.l.Lock()
	c := q.queue.Len()
	q.l.Unlock()

	return c
}

// Add an item to the end of the queue
func (q *LockingQueue) Add(item interface{}) {
	// locking to make Queue thread-safe
	q.l.Lock()
	q.queue.Add(item)
	q.l.Unlock()
}

func (q *LockingQueue) AddList(items []interface{}) {
	q.l.Lock()
	q.queue.AddList(items)
	q.l.Unlock()
}

// Remove the item at the head of the queue and return it.
// Returns nil when there are no items left in queue.
func (q *LockingQueue) Next() (item interface{}) {
	// locking to make Queue thread-safe
	q.l.Lock()
	i := q.queue.Next()
	q.l.Unlock()

	return i
}

// Returns the next N elements from the queue
// In case of not enough elements, returns the elements that are available
func (q *LockingQueue) NextN(n int) []interface{} {
	q.l.Lock()
	i := q.queue.NextN(n)
	q.l.Unlock()

	return i
}

// Reads the item at the head of the queue without removing it
// Returns nil when there are no items left in queue
func (q *LockingQueue) Peek() (item interface{}) {
	q.l.Lock()
	i := q.queue.Peek()
	q.l.Unlock()

	return i
}
