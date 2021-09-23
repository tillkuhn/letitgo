package main

import (
	"log"
	"sync"
)

func main() {
	simulateRaceCondition()
}

// https://github.com/MicahParks/keyfunc/issues/16
// Some people are familiar with the map read/write race condition. For those who are not,
//  take a look at this example. It should cause a Go panic.
func simulateRaceCondition() {
	m := make(map[string]int)
	wait := sync.WaitGroup{}
	iterations := 1000
	const key = ""
	wait.Add(1)
	go func() {
		for i := 0; i < iterations; i++ {
			m[key] = i
		}
		wait.Done()
	}()
	for i := 0; i < iterations; i++ {
		log.Printf("key: %d", m[key])
	}
	wait.Wait()
}
