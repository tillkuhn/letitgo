package cache

import (
	"testing"
	"time"
)

func TestCacheLoop(t *testing.T) {
	cacheLoop(250*time.Millisecond, 100*time.Millisecond, 4)
}
