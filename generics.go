// Created by Yaz Saito on 06/15/12.
// Modified by Geert-Johan Riemer, Foize B.V., Ferenc Fabian

// TODO:
// - travis CI
// - maybe add method (*Queue).Peek()

package fifo

import (
	"reflect"
	"sync"
)

// genericChunk are used to make a queue auto resizeable.
type genericChunk[T any] struct {
	items       [chunkSize]T     // list of queue'ed items
	first, last int              // positions for the first and list item in this chunk
	next        *genericChunk[T] // pointer to the next chunk (if any)
}

// GenericQueue is an implementation of fifo queue with generic parameters
type GenericQueue[T any] struct {
	head, tail *genericChunk[T] // chunk head and tail
	count      int              // total amount of items in the queue
	lock       sync.Mutex       // synchronisation lock
}

// NewGenericQueue creates a new and empty *fifo.GenericQueue[T]
func NewGenericQueue[T any]() (q *GenericQueue[T]) {
	initChunk := new(genericChunk[T])
	q = &GenericQueue[T]{
		head: initChunk,
		tail: initChunk,
	}
	return q
}

// Len returns the number of items in the queue
func (q *GenericQueue[T]) Len() (length int) {
	// locking to make Queue thread-safe
	q.lock.Lock()
	defer q.lock.Unlock()

	// copy q.count and return length
	length = q.count
	return length
}

// Add an item to the end of the queue
func (q *GenericQueue[T]) Add(item T) {
	// locking to make Queue thread-safe
	q.lock.Lock()
	defer q.lock.Unlock()

	// check if item is valid
	if reflect.ValueOf(&item).Elem().IsZero() {
		panic("can not add nil item to fifo queue")
	}

	// if the tail chunk is full, create a new one and add it to the queue.
	if q.tail.last >= chunkSize {
		q.tail.next = new(genericChunk[T])
		q.tail = q.tail.next
	}

	// add item to the tail chunk at the last position
	q.tail.items[q.tail.last] = item
	q.tail.last++
	q.count++
}

// Next removes the item at the head of the queue and return it.
// Returns nil when there are no items left in queue.
func (q *GenericQueue[T]) Next() (item T) {
	// locking to make Queue thread-safe
	q.lock.Lock()
	defer q.lock.Unlock()

	// Return nil if there are no items to return
	if q.count == 0 {
		var empty T
		return empty
	}
	// FIXME: why would this check be required?
	if q.head.first >= q.head.last {
		var empty T
		return empty
	}

	// Get item from queue
	item = q.head.items[q.head.first]

	// increment first position and decrement queue item count
	q.head.first++
	q.count--

	if q.head.first >= q.head.last {
		// we're at the end of this chunk and we should do some maintainance
		// if there are no follow up chunks then reset the current one so it can be used again.
		if q.count == 0 {
			q.head.first = 0
			q.head.last = 0
			q.head.next = nil
		} else {
			// set queue's head chunk to the next chunk
			// old head will fall out of scope and be GC-ed
			q.head = q.head.next
		}
	}

	// return the retrieved item
	return item
}
