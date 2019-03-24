// Created by Yaz Saito on 06/15/12.
// Modified by Geert-Johan Riemer, Foize B.V.

// TODO:
// - travis CI

package fifo

const chunkSize = 64

// chunks are used to make a queue auto resizeable.
type chunk struct {
	items       [chunkSize]interface{} // list of queue'ed items
	first, last int                    // positions for the first and list item in this chunk
	next        *chunk                 // pointer to the next chunk (if any)
}

// fifo queue
type UnsafeQueue struct {
	head, tail *chunk // chunk head and tail
	count      int    // total amount of items in the queue
}

// NewUnsafeQueue creates a new and empty *fifo.UnsafeQueue
func NewUnsafeQueue() (q *UnsafeQueue) {
	initChunk := new(chunk)
	q = &UnsafeQueue{
		head: initChunk,
		tail: initChunk,
	}
	return q
}

// Return the number of items in the queue
func (q *UnsafeQueue) Len() (length int) {
	// copy q.count and return length
	length = q.count
	return length
}

// Add an item to the end of the queue
func (q *UnsafeQueue) Add(item interface{}) {
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

// Adds an list of items to the queue
func (q *UnsafeQueue) AddList(items []interface{}) {
	// check if item is valid
	if len(items) == 0 {
		// len(nil) == 0 as well
		return
	}
	//
	if len(items) > chunkSize { // Add each piece separated
		chunks := len(items) / chunkSize
		if chunks*chunkSize != len(items) { // Rouding up
			chunks++
		}

		for i := 0; i < chunks; i++ {
			s := i * chunkSize
			e := (i + 1) * chunkSize

			if e > len(items) {
				e = len(items)
			}

			q.AddList(items[s:e])
		}
		return
	}

	// if the tail chunk is full, create a new one and add it to the queue.
	if q.tail.last >= chunkSize {
		q.tail.next = new(chunk)
		q.tail = q.tail.next
	}

	s := q.tail.last
	e := len(items) - s
	n := copy(q.tail.items[s:e], items)
	q.tail.last += n
	q.count += n
	items = items[e:]

	if len(items) > 0 {
		q.AddList(items) // Add Remaining Items
	}
}

// Returns the next N elements from the queue
// In case of not enough elements, returns the elements that are available
func (q *UnsafeQueue) NextN(n int) []interface{} {
	if n > chunkSize {
		// Recursive call to append
		chunks := n / chunkSize
		if chunks*chunkSize < n {
			chunks++
		}

		out := make([]interface{}, 0)
		read := 0
		for i := 0; i < chunks; i++ {
			e := chunkSize
			if read+e > n {
				e = n - read
			}

			out = append(out, q.NextN(e)...)
		}
		return out
	}

	if q.count < n {
		n = q.count // Not enough elements
	}

	if q.count == 0 || q.head.first >= q.head.last {
		return make([]interface{}, 0)
	}

	// TODO: Slice it
	out := make([]interface{}, n)

	read := 0
	for i := 0; i < n; i++ {
		if q.count == 0 {
			break
		}
		read++
		out[i] = q.Next()
	}

	return out[:read]
}

// Remove the item at the head of the queue and return it.
// Returns nil when there are no items left in queue.
func (q *UnsafeQueue) Next() (item interface{}) {

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

// Reads the item at the head of the queue without removing it
// Returns nil when there are no items left in queue
func (q *UnsafeQueue) Peek() (item interface{}) {
	// Return nil if there are no items to return
	if q.count == 0 {
		return nil
	}
	// FIXME: why would this check be required?
	if q.head.first >= q.head.last {
		return nil
	}

	// Get item from queue
	return q.head.items[q.head.first]
}
