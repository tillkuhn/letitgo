package main

import (
	"fmt"
	"sync"
)

// Stack is a LIFO stack that does not grow bigger than maxLen
type Stack struct {
	mu sync.Mutex
	items []int
	maxLen int
}

// Add element to the end of the stack, shrink if necessary
func (s *Stack) add(i int) {
	// For is Go's "while"
	s.mu.Lock()
	defer s.mu.Unlock()
	for len(s.items) >= s.maxLen {
		first := s.items[0]
		copy(s.items, s.items[1:])
		s.items = s.items[:len(s.items)-1]
		fmt.Printf("%d popped\n", first)
	}
	s.items = append(s.items, i)
	fmt.Printf("%v newlen %d\n", s.items, len(s.items))
}

