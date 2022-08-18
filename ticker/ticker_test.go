package ticker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTicker(t *testing.T) {
	var tickCount int
	tickerFunc := func(t time.Time) {
		tickCount++
	}
	doneChan := RunTickerFuncWithChannel(tickerFunc, 5*time.Millisecond)
	time.Sleep(50 * time.Millisecond)
	doneChan <- true
	assert.Greater(t, tickCount, 0)
}
