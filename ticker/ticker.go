package ticker

import (
	"fmt"
	"time"
)

// RunTickerWithChannel ticker example based on https://gobyexample.com/tickers
// Tickers use a similar mechanism to timers: a channel that is sent values.
// Here we’ll use the select builtin on the channel to await the values as they arrive every 500ms.
func RunTickerWithChannel() {

	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()

	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}

// RunTickerSimple simplified version based on https://stackoverflow.com/a/53057336/4292075
func RunTickerSimple() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			fmt.Println("Hello !!")
		}
	}()

	// wait for 10 seconds
	time.Sleep(10 * time.Second)
	ticker.Stop()
}
