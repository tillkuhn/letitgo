// Package concurrent: JOB QUEUES IN GO
// Constructs and snippets to build your job queue in Golang.
// Good Article: https://www.opsdash.com/blog/job-queues-in-go.htm
package concurrent

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

type Job struct {
	Id  int
	Seq int
}

func (j Job) String() string {
	return fmt.Sprintf("Seq #%d ID=%d", j.Seq, j.Id)
}

func worker(jobChan <-chan Job) {
	defer wg.Done()
	fmt.Println("Worker: Start worker")
	for job := range jobChan {
		fmt.Printf("Worker: Processing Job %v\n", job)
		// artificial time while we do "work"
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("Worker: No more work")
}

// DoWork employs a Worker with a buffered Job Queue, see
// https://www.opsdash.com/blog/job-queues-in-go.html
func DoWork() {
	// make a channel with a capacity of 10.
	var capacity = 10
	jobChan := make(chan Job, capacity)

	// start the worker but
	// increment the WaitGroup before starting the worker
	fmt.Println("Chef: Start worker")
	wg.Add(1)
	go worker(jobChan)

	// enqueue a couple jobs, make sure we have more than capacity
	for i, id := range []int{1, 32, 45, 667, 22, 444, 6636, 77, 99, 11, 101, 102, 103} {
		fmt.Printf("Chef: Enqueue id %d\n", id)
		job := Job{Seq: i, Id: id}
		// jobChan <- job
		// ENQUEUEING WITHOUT BLOCKING
		if res := tryEnqueue(job, jobChan); !res {
			fmt.Printf("Chef: Job capcatity %d exceeded, job %v could not be enqueued\n", capacity, job)
		}
	}

	// to stop the worker, first close the job channel
	fmt.Println("Chef: Stop worker")
	close(jobChan)

	// then wait using the WaitGroup
	wg.Wait()
	fmt.Println("Chef: Time to call it a day")

}

// tryEnqueue tries to enqueue a job to the given job channel. Returns true if
// the operation was successful, and false if enqueuing would not have been
// possible without blocking. Job is not enqueued in the latter case.
func tryEnqueue(job Job, jobChan chan<- Job) bool {
	select {
	case jobChan <- job:
		return true
	default:
		return false
	}
}
