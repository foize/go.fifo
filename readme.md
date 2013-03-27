## go.fifo

### Description
go.fifo provides a simple FIFO thread-safe queue.
*fifo.Queue supports pushing an item at the end with Add(), and popping an item from the front with Next().
There is no intermediate type for the stored data. Data is directly added and retrieved as type interface{}
The queue itself is implemented as a single-linked list of chunks containing max 64 items each.

### Installation
`go get github.com/foize/go.fifo`

### Usage
```go
package main

import (
	"github.com/foize/go.fifo"
	"fmt"
)

func main() {
	// create a new queue
	queue := fifo.NewQueue()

	// add items to the queue
	queue.Add(42)
	queue.Add(123)
	queue.Add(456)

	// retrieve items from the queue
	fmt.Println(queue.Next().(int)) // 42
	fmt.Println(queue.Next().(int)) // 123
	fmt.Println(queue.Next().(int)) // 456
}
```

### Documentation
Documentation can be found at [godoc.org/github.com/foize/go.fifo](http://godoc.org/github.com/foize/go.fifo).
For more detailed documentation, read the source.

### History
This package is based on github.com/yasushi-saito/fifo_queue
There are several differences:
- renamed package to `fifo` to make usage simpler
- removed intermediate type `Item` and now directly using interface{} instead.
- renamed (*Queue).PushBack() to (*Queue).Add()
- renamed (*Queue).PopFront() to (*Queue).Next()
- Next() will not panic on empty queue, will just return nil interface{}
- Add() does not accept nil interface{} and will panic when trying to add nil interface{}.
- Made fifo.Queue thread/goroutine-safe (sync.Mutex)
- Added a lot of comments
- renamed internal variable/field names