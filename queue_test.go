package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// https://groups.google.com/g/golang-nuts/c/obZI4uyZTe0
// https://golang.org/src/database/sql/sql.go#L791
func TestStack(t *testing.T) {
	s := Queue{
		items:  intSlice2Interface([]int{1, 2, 4, 5, 6, 7, 8, 999, 9993}),
		maxLen: 6,
	}
	s.Enqueue(888)
	assert.GreaterOrEqualf(t, s.maxLen, len(s.items), "not > maxLen")
	s.Enqueue(8884)
	s.Enqueue(88555)
	s.Enqueue(899995)
	assert.GreaterOrEqualf(t, s.maxLen, len(s.items), "not > maxLen")
	p := s.Dequeue().(int)
	assert.Equal(t, 999, p)
}

// https://stackoverflow.com/a/12754757/4292075
func intSlice2Interface(a []int) []interface{} {
	b := make([]interface{}, len(a))
	for i := range a {
		b[i] = a[i]
	}
	return b
}
