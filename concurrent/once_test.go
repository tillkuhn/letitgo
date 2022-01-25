package concurrent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMainFunc(t *testing.T) {
	DoItOnceSam()
	DoItOnceSam()
	assert.Equal(t, 42, initNumber, "number should've been incremented exactly once")

}
