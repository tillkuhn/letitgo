package concurrent

import (
	"log"
	"runtime"
	"time"
)

// WaitAndExec https://www.geeksforgeeks.org/time-afterfunc-function-in-golang-with-examples/
func WaitAndExec() {
	// Defining duration parameter of
	// AfterFunc() method
	callAfter := time.Duration(1) * time.Second

	// Defining function parameter of
	// AfterFunc() method
	f := func() {
		// Printed when it's called by the
		// AfterFunc() method in the time
		// stated above
		log.Printf("Function called by "+
			"AfterFunc() after %v seconds  goroutines=%d\n", callAfter, runtime.NumGoroutine())
	}

	// Calling AfterFunc() method with its
	// parameter
	log.Printf("Setting Timer for %v goroutines=%d", callAfter, runtime.NumGoroutine())
	timer1 := time.AfterFunc(callAfter, f)

	// Calling stop method
	// w.r.to timer1
	defer timer1.Stop()

	// Calling sleep method
	sleepy := callAfter + 1*time.Second
	log.Printf("Sleeping %v\n", sleepy)
	time.Sleep(sleepy)
	log.Printf("Finished, goroutines=%d\n", runtime.NumGoroutine())
}
