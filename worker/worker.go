// Package worker first attempt to "do something with go generics" :-)
package worker

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
)

// Worker manages all bits and pieces
// T represents the Job Type, i.e. anything that can be processed by JobFunc,
// but should implement Stringer Interface for compact logging
//
// Example Usage:
//
//	 w := Worker[Job]{}
//	 w.Start(func(job Job){ fmt.Printf("Processing Job %v",job) },10)
//	 for i, id := range []int{1, 32, 45, 667, 22} {
//	    w.TryEnqueue(Job{Seq: i, Id: id})
//	  }
//	  w.Stop()
//
type Worker[T any] struct {
	jobFunc func(T) error
	jobChan chan T
	wg      sync.WaitGroup
	logger  zerolog.Logger
}

// Start employs a Worker goroutine attached to a buffered Job Queue
func (w *Worker[T]) Start(jobFunc func(T) error, capacity int) {
	w.jobChan = make(chan T, capacity)
	w.jobFunc = jobFunc
	w.logger = log.With().Str("logger", "worker").Logger()
	w.logger.Debug().Msgf("Chef: Start worker with channel capacity %d", capacity)
	w.wg.Add(1) // increment for each worker, currently we keep it simple and only employ a single worker
	go w.listen()
}

// Stop closes the worker channel and waits for the WaitGroup
func (w *Worker[T]) Stop() {
	// to stop the worker, first close the job channel
	w.logger.Debug().Msg("Chef: Stop worker")
	close(w.jobChan)
	w.wg.Wait()
	w.logger.Debug().Msg("Chef: Time to call it a day")
}

// QueueLength returns info on the current length of the job channel (backlog)
// IntelliJ reports "Invalid argument for the len function" for generic args, which is a false positive (same with cap)
func (w *Worker[T]) QueueLength() int {
	return len(w.jobChan)
}

// Capacity returns info on the current capacity of the job channel
func (w *Worker[T]) Capacity() int {
	return cap(w.jobChan)
}

// TryEnqueue tries to enqueue a job to the given job channel.
// Returns "true" if the operation was successful
// Returns "false" if enqueuing would not have been possible without blocking.
// Job is not enqueued in the latter case.
func (w *Worker[T]) TryEnqueue(job T) bool {
	select {
	case w.jobChan <- job:
		return true
	default:
		w.logger.Debug().Msgf("jobChan is full, cannot enqueue %v", job)
		return false
	}
}

// listen sits on the job channel, and executes the jobFunction for each incoming job
// errors are currently only logged, but not reported to a dedicated channel
func (w *Worker[T]) listen( /*jobChan <-chan*/ ) {
	logger := log.With().Str("logger", "worker").Logger()
	defer w.wg.Done()
	logger.Debug().Msg("Worker: Start listening")
	for job := range w.jobChan {
		logger.Debug().Msgf("Worker: Processing Job %v", job)
		if err := w.jobFunc(job); err != nil {
			logger.Error().Msgf("Error during job processing: %v ", err)
		}
		logger.Debug().Msgf("Worker: Processed Job %v", job)
	}
	logger.Debug().Msg("Worker: No more work")
}
