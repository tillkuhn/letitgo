package concurrent

import (
	"sync"
)

var (
	once       sync.Once // guards initMime
	initNumber int
)

// DoItOnceSam should call initStuff only once
func DoItOnceSam() {
	once.Do(initStuff)
	// do something else
}

func initStuff() {
	initNumber = initNumber + 42
}
