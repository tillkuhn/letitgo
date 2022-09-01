package filesystem

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoin(t *testing.T) {
	assert.Equal(t, "hase/horst", join("hase", "horst"))
}
