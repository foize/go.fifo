// Created by Yaz Saito on 06/15/12.
// Modified by Geert-Johan Riemer, Foize B.V.

// go.fifo provides a simple FIFO thread-safe queue.
// *fifo.Queue supports pushing an item at the end with Add(), and popping an item from the front with Next().
// There is no intermediate type for the stored data. Data is directly added and retrieved as type interface{}
// The queue itself is implemented as a single-linked list of chunks containing max 64 items each.

// This package is based on github.com/yasushi-saito/fifo_queue
// There are several differences:
// - renamed package to `fifo` for easynes of use
// - removed intermediate type `Item` and now directly using interface{} instead.
// - renamed (*Queue).PushBack() to (*Queue).Add()
// - renamed (*Queue).PopFront() to (*Queue).Next()
// - Next() will not panic on empty queue, will just return nil interface{}
// - Add() does not accept nil interface{} and will panic when trying to add nil interface{}.
// - Made fifo.Queue thread/goroutine-safe (sync.Mutex)
// - Added a lot of comments
// - renamed internal variable/field names

// TODO:
// - travis CI
// - maybe add method (*Queue).Peek()

package fifo

import (
	"sync"
)

const chunkSize = 64

// chunks are used to make a queue auto resizeable.
type chunk struct {
	items       [chunkSize]interface{} // list of queue'ed items
	first, last int                    // positions for the first and list item in this chunk
	next        *chunk                 // pointer to the next chunk (if any)
}

// fifo queue
type Queue struct {
	head, tail *chunk     // chunk head and tail
	count      int        // total amount of items in the queue
	lock       sync.Mutex // synchronisation lock
}

// NewQueue creates a new and empty *fifo.Queue
func NewQueue() (q *Queue) {
	initChunk := new(chunk)
	q = &Queue{
		head: initChunk,
		tail: initChunk,
	}
	return q
}

// Return the number of items in the queue
func (q *Queue) Len() (length int) {
	// locking to make Queue thread-safe
	q.lock.Lock()
	defer q.lock.Unlock()

	// copy q.count and return length
	length = q.count
	return length
}

// Add an item to the end of the queue
func (q *Queue) Add(item interface{}) {
	// locking to make Queue thread-safe
	q.lock.Lock()
	defer q.lock.Unlock()

	// check if item is valid
	if item == nil {
		panic("can not add nil item to fifo queue")
	}

	// if the tail chunk is full, create a new one and add it to the queue.
	if q.tail.last >= chunkSize {
		q.tail.next = new(chunk)
		q.tail = q.tail.next
	}

	// add item to the tail chunk at the last position
	q.tail.items[q.tail.last] = item
	q.tail.last++
	q.count++
}

// Remove the item at the head of the queue and return it.
// Returns nil when there are no items left in queue.
func (q *Queue) Next() (item interface{}) {
	// locking to make Queue thread-safe
	q.lock.Lock()
	defer q.lock.Unlock()

	// Return nil if there are no items to return
	if q.count == 0 {
		return nil
	}
	// FIXME: why would this check be required?
	if q.head.first >= q.head.last {
		return nil
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
