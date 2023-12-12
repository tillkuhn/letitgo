package worker

import (
	"fmt"
)

type MyJob struct {
	Id  int
	Seq int
}

func DoWork() {
	ids := []int{1, 32, 45, 667, 22}
	w := Worker[MyJob]{}
	w.Start(func(job MyJob) error {
		fmt.Printf("Just doing my job %v\n", job)
		return nil
	}, len(ids))
	for i, id := range ids {
		w.TryEnqueue(MyJob{Seq: i, Id: id})
	}
	w.Stop()
}
