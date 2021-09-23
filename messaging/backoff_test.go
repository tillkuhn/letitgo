package messaging

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBackoff(t *testing.T) {
	err := RunWithBackoff(3)
	assert.Error(t, err)
}
