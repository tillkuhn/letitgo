package messaging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBackoff(t *testing.T) {
	err := RunWithBackoff(3)
	assert.Error(t, err)
}
