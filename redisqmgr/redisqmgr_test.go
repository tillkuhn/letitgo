package redisqmgr

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

type ScanTask struct {
	ScanId int `json:"scanId"`
}

func (st ScanTask) MarshalBinary() (data []byte, err error) {
	return json.Marshal(st)
}

// TestRedisQueueRealInstance skipped in ci, expects a real redis instance
// Todo check out https://github.com/go-redis/redismock
func TestRedisQueueRealInstance(t *testing.T) {
	if _, present := os.LookupEnv("CI_PROJECT_NAME"); present {
		t.Skip("Skipping testing in CI environment")
	}
	client, mock := redismock.NewClientMock()
	rq := NewWithClient(client, &ClientOpts{Namespace: "Test"})
	qShort := "kuh-queue"
	qFull := rq.QueueWithNamespace(qShort)
	task1 := ScanTask{ScanId: 14928080}
	task2 := ScanTask{ScanId: 47118080}

	// mock.ExpectGet(key).RedisNil()
	// 	mock.Regexp().ExpectSet(key, `[a-z]+`, 30 * time.Minute).SetErr(errors.New("FAIL"))
	task1m, err := task1.MarshalBinary()
	assert.NoError(t, err)
	task2m, err := task2.MarshalBinary()
	assert.NoError(t, err)

	mock.ExpectRPush(qFull, string(task1m)).SetVal(1)
	mock.ExpectRPush(qFull, string(task2m)).SetVal(2) // returns list length after push
	mock.ExpectBLPop(rq.listenerTimeout, qFull).SetVal([]string{qShort, string(task2m)})
	mock.ExpectBLPop(rq.listenerTimeout, qFull).SetVal([]string{qShort, string(task1m)})

	ctx, cancel := context.WithCancel(context.Background())
	err = rq.Push(ctx, qShort, task1)
	assert.NoError(t, err)
	assert.NoError(t, rq.Push(ctx, qShort, task2))

	go rq.StartQueueWorker(ctx, qShort, func(s string) error {
		assert.Contains(t, s, "8080") // all IDs contain 8080
		expectErr := mock.ExpectationsWereMet()
		if expectErr == nil {
			cancel() // all messages processed, no need to continue polling during test
		}
		return nil
	}, 100*time.Millisecond)
	time.Sleep(1 * time.Second)
	cancel()                     // Signal cancellation to context.Context
	rq.WaitForListenerShutdown() // Block here until are workers are done
	log.Info().Msgf("Wg done")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestQueueNameQualifiedAndWithoutNamespace(t *testing.T) {
	rq := New(&ClientOpts{Namespace: "testa"})
	assert.Equal(t, "testa.hase-klaus", rq.QueueWithNamespace("hase-klaus"))
	rq = New(&ClientOpts{Namespace: ""})
	assert.Equal(t, "hase-horst", rq.QueueWithNamespace("hase-horst"))
}

//func TestBlockingWithCancel(t *testing.T) {
//	ctx, cancel := context.WithCancel(context.Background())
//	// start goroutine with a new context
//	go handleRequest(ctx, 1)
//	time.Sleep(2 * time.Second) // Time between requests
//	cancel()
//}
//
//func handleRequest(ctx context.Context, incr int) {
//	fmt.Println("New request registered: ", incr+1)
//	for i := 0; i <= 100; i++ {
//		fmt.Println("Request: ", incr+1, " | Sub-task: ", i+1)
//		time.Sleep(100 * time.Millisecond) // Time processing
//		select {
//		case <-ctx.Done():
//			// canceled
//			fmt.Println("Canceled")
//			return
//		default:
//			fmt.Println("Not Canceled")
//		}
//	}
//	return
//}
