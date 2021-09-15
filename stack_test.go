package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// https://groups.google.com/g/golang-nuts/c/obZI4uyZTe0
// https://golang.org/src/database/sql/sql.go#L791
func TestStack(t *testing.T) {
	s := Stack{ items: []int{1, 2,4,5,6,7,8, 999, 9993}, maxLen: 6}
	s.add(88)
	assert.GreaterOrEqualf(t, s.maxLen, len(s.items),"not > maxLen")
	s.add(884)
	s.add(88555)
	s.add(899995)
	assert.GreaterOrEqualf(t, s.maxLen, len(s.items),"not > maxLen")

}
