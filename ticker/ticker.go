package ticker

import (
	"fmt"
	"time"
)

// RunTickerWithChannel ticker example based on https://gobyexample.com/tickers
// Tickers use a similar mechanism to timers: a channel that is sent values.
// Here we’ll use the select builtin on the channel to await the values as they arrive every 500ms.
// https://stackoverflow.com/a/31936175/4292075
func RunTickerWithChannel() {
	var tickCount int
	tickerFunc := func(t time.Time) {
		tickCount++
		fmt.Println("Tick tock at", t, "tickCount=", tickCount)
	}
	doneChan := RunTickerFuncWithChannel(tickerFunc, 500*time.Millisecond)
	time.Sleep(1600 * time.Millisecond)
	doneChan <- true
	fmt.Printf("Function has been executed %d times\n", tickCount)
	time.Sleep(1600 * time.Millisecond)
}

// RunTickerFuncWithChannel runs the ticket function once immediately, afterwards in regular intervals
// until
func RunTickerFuncWithChannel(tickerFunc func(t time.Time), interval time.Duration) chan bool {
	doneChan := make(chan bool)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		tickerFunc(time.Now())
		for {
			select {
			case <-doneChan:
				fmt.Println("We're done here")
				return
			case t := <-ticker.C:
				tickerFunc(t)
			}
		}
	}()
	return doneChan
}

// RunTickerSimple simplified version based on https://stackoverflow.com/a/53057336/4292075
// func RunTickerSimple() {
//	ticker := time.NewTicker(1 * time.Second)
//	go func() {
//		for range ticker.C {
//			fmt.Println("Hello !!")
//		}
//	}()
//
//	// wait for 10 seconds
//	time.Sleep(10 * time.Second)
//	ticker.Stop()
//}
