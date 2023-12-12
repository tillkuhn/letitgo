// Package worker first attempt to "do something with go generics" :-)
package worker

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type PizzaJob struct {
	Id      int
	Seq     int
	Topping string
}

var ids = []int{1, 32, 45, 667, 22}

func TestWorker(t *testing.T) {
	w := NewWorker[PizzaJob]()
	w.Start(justDoingMyJob, len(ids))
	for i, id := range ids {
		res := w.TryEnqueue(PizzaJob{Seq: i, Id: id})
		assert.True(t, res, " PizzaJob capacity should be sufficient")
	}
	assert.Equal(t, len(ids), w.Capacity())
	w.Stop()
}

func TestWorkerBlocking(t *testing.T) {
	w := Worker[PizzaJob]{}
	unblock := false
	w.Start(justDoingMyBlockingJob(&unblock), len(ids))
	for i, id := range ids {
		res := w.TryEnqueue(PizzaJob{Seq: i, Id: id})
		assert.True(t, res, " PizzaJob capacity should be sufficient")
	}
	// add another one - this should be enough to exceed the queue size
	assert.Equal(t, len(ids), w.QueueLength())
	res := w.TryEnqueue(PizzaJob{Seq: 999, Id: 99})
	assert.False(t, res)
	assert.Equal(t, len(ids), w.QueueLength()) // should be still the same

	unblock = true // release the waiting jobs
	w.Stop()
}

func justDoingMyJob(job PizzaJob) error {
	fmt.Printf("Jobber %v", job)
	if job.Id < 1 {
		return errors.New("I don't like your Id")
	}
	return nil
}

// justDoingMyBlockingJob takes boolean which - unless true - blocks the jon execution
// this allows us to make predictable tests on queue size, and release the jobs subsequently
func justDoingMyBlockingJob(unblock *bool) func(job PizzaJob) error {
	return func(job PizzaJob) error {
		for {
			if *unblock {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		return nil
	}
}
