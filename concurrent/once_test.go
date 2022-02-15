package concurrent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainFunc(t *testing.T) {
	DoItOnceSam()
	DoItOnceSam()
	assert.Equal(t, 42, initNumber, "number should've been incremented exactly once")

}
