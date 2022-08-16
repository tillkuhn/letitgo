package ticker

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
