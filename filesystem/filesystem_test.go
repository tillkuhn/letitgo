package filesystem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	assert.Equal(t, "hase/horst", join("hase", "horst"))
}
