package messaging

import (
	"fmt"
	"log"
	"sync"
)

// Queue A queue is a linear data structure in which elements can be inserted only
// from one side of the list called rear, and the elements can be deleted only from
// the other side called the front. The queue data structure follows the FIFO (First In First Out) principle,
// i.e. the element inserted at first in the list, is the first element to be removed from the list.
// The insertion of an element in a queue is called an enqueue operation and the deletion of an element
// is called a dequeue operation.
// Source: https://www.geeksforgeeks.org/difference-between-stack-and-queue-data-structures/
//
// This is generic (interface{}) but you can simply replace the type of the slice elements
// with something more concrete
type Queue struct {
	mu     sync.Mutex
	items  []interface{}
	maxLen int
}

// Enqueue Add element to the end of the stack, shrink if necessary
func (s *Queue) Enqueue(i interface{}) {
	// For is Go's "while"
	for len(s.items) >= s.maxLen {
		popped := s.Dequeue()
		log.Printf("Element %v was dequed", popped)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, i)
	fmt.Printf("%v newlen %d\n", s.items, len(s.items))
}

func (s *Queue) Dequeue() interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	first := s.items[0]
	// copy(s.items, s.items[1:])
	// s.items = s.items[:len(s.items)-1]
	s.items = s.items[1:] // use simpler re-slice approach
	return first

}
