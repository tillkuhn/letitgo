package redisqmgr

import (
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"

	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// QueueMgr messaging queueing backed by https://github.com/go-redis/redis inspired by ...
// "Implementing Message Queuing with Golang, Redis, and MySQL 8 on Linux Server"
// https://www.vultr.com/docs/implementing-message-queuing-with-golang-redis-and-mysql-8-on-linux-server/
// "Golang Simple Job Queue With Redis Streams"
// https://codesahara.com/blog/golang-job-queue-with-redis-streams/
// More Related Topics
// https://ably.com/blog/event-streaming-with-redis-and-golang
// https://callistaenterprise.se/blogg/teknik/2019/10/05/go-worker-cancellation/
type QueueMgr struct {
	redisClient     *redis.Client
	logger          zerolog.Logger
	listenerWG      sync.WaitGroup // no need to initialize
	listenerTimeout time.Duration
	gaugeQueued     *prometheus.GaugeVec
	counterError    *prometheus.CounterVec
	counterSuccess  *prometheus.CounterVec
	clientOpts      *ClientOpts
}

// ClientOpts  bundles the options for creating Queue Manager Clients
type ClientOpts struct {
	RedisAddr     string
	RedisPassword string
	Namespace     string
}

// New returns a new Queue Manager based on default redis address on localhost
// delegates to NewWithClient
func New(opts *ClientOpts) *QueueMgr {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     opts.RedisAddr,
		Password: opts.RedisPassword,
		DB:       0,
	})
	return NewWithClient(redisClient, opts)
}

// NewWithClient returns a new Queue Manager, bring-your-own-redisClient style
// also registers a prometheus gauge vec as shown in
// https://github.com/prometheus/client_golang/blob/main/prometheus/examples_test.go#L51
func NewWithClient(client *redis.Client, opts *ClientOpts) *QueueMgr {
	// Init prometheus Gauges and Counters
	gaugeQueued := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: opts.Namespace,
			Name:      "ops_queued",
			Help:      "Number of operations waiting to be processed, partitioned by queue name",
		},
		[]string{"queue"},
	)
	// CounterVec open issue https://github.com/prometheus/client_golang/issues/190
	// so unless there is at least one error added, the counter for that label won't show in metrics export
	counterError := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: opts.Namespace,
			Name:      "ops_error",
			Help:      "Number of operations than run into an error, partitioned by queue name",
		},
		[]string{"queue"},
	)
	counterSuccess := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: opts.Namespace,
			Name:      "ops_success",
			Help:      "Number of operations than have been successfully processed, partitioned by queue name",
		},
		[]string{"queue"},
	)
	prometheus.MustRegister(gaugeQueued)
	prometheus.MustRegister(counterError)
	prometheus.MustRegister(counterSuccess)

	rq := QueueMgr{
		redisClient:     client,
		logger:          log.With().Str("logger", "redis-qmgr").Logger(),
		listenerTimeout: 1 * time.Second,
		gaugeQueued:     gaugeQueued,
		counterError:    counterError,
		counterSuccess:  counterSuccess,
		clientOpts:      opts,
	}
	return &rq
}

// Push pushes data to the given "queue" synchronously
// This is backed by a redis list value, see docs https://redis.io/commands/rpush/
// RPush: "Insert all the specified values at the tail of the list stored at key.
// If key does not exist, it is created as empty list before performing the push operation."
// Events must implement binary Marshaller interface, so they can be serialized properly
// RPush returns the length of the list after the push operation.
func (rq *QueueMgr) Push(ctx context.Context, queueShort string, event encoding.BinaryMarshaler) error {
	queue := rq.QueueWithNamespace(queueShort)
	rq.logger.Info().Msgf("Pushing message to queue=%s", queue)
	b, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if err := rq.redisClient.RPush(ctx, queue, string(b)).Err(); err != nil {
		rq.counterError.WithLabelValues(queueShort).Add(1)
		return err
	}
	rq.logger.Info().Msgf("Push successful, queue=%s, qLen=%d", queue, rq.QueueLen(ctx, queueShort))
	return nil
}

// StartQueueWorker listen on a queue represented by a Redis List Value
// This method starts a go routine and returns immediately
// It uses BLPop internally to read "Messages" with a small timeout, see https://redis.io/commands/blpop/
// "BLPOP is a blocking list pop primitive. It is the blocking version of LPOP
// because it blocks the connection when there are no elements to pop from any of the given lists
func (rq *QueueMgr) StartQueueWorker(ctx context.Context, queueShort string, callbackFunc func(string) error, pollInterval time.Duration) {
	rq.listenerWG.Add(1)
	go rq.queueWorker(ctx, queueShort, callbackFunc, pollInterval)
}

// queueWorker performs the actual, started by StartQueueWorker in a separate goroutine
func (rq *QueueMgr) queueWorker(ctx context.Context, queueShort string, callbackFunc func(string) error, pollInterval time.Duration) {
	// decrement workgroup listener counter by one once this worker exits
	defer rq.listenerWG.Done()
	queue := rq.QueueWithNamespace(queueShort)
	logger := log.With().Str("logger", "redis-"+queueShort).Logger()
	logger.Info().Msgf("Start polling messages from queue=%s queued=%d interval=%s",
		queue, rq.QueueLen(ctx, queueShort), pollInterval)
	for {
		// https://stackoverflow.com/questions/52279988/how-to-do-redis-redigo-lpop-in-golang
		// command returns a two element array where the first element is the key and the second value is the popped element.
		results, err := rq.redisClient.BLPop(ctx, rq.listenerTimeout, queue).Result()
		// redis.nil is OK if the key doesn't exist
		if err != nil && err != redis.Nil {
			logger.Error().Msgf("Error during pop %s: %v", queue, err)
		}
		if len(results) > 1 { // expect 2, fist is the queue name, 2nd is the value
			qLen := rq.QueueLen(ctx, queueShort)
			logger.Debug().Msgf("Processing msg (total=%d) from %s: %d bytes", qLen, queue, len(results[1]))
			if err := callbackFunc(results[1]); err != nil {
				rq.counterError.WithLabelValues(queueShort).Add(1)
				logger.Error().Msgf("error during job execution: %v", err)
			} else {
				rq.counterSuccess.WithLabelValues(queueShort).Add(1)
			}
		}
		select {
		case <-ctx.Done():
			// canceled
			logger.Info().Msgf("Context canceled, disconnecting from queue=%s", queue)
			return
		default:
			logger.Trace().Msgf("Not cancelled, sleep %ds before popping again", pollInterval/1000)
			time.Sleep(pollInterval)
		}
	}
}

// QueueLen Returns the number of list elements in the key identified by queue name, or zero if the key does not exist
// this method has no side effects, except updating the Counter Gauge for queue length
// Since the Counter is already namespaced, we only use the short queue name!
func (rq *QueueMgr) QueueLen(ctx context.Context, queueShort string) int {
	qLen := rq.redisClient.LLen(ctx, rq.QueueWithNamespace(queueShort)).Val()
	// Update Gauge with current length.
	rq.gaugeQueued.WithLabelValues(queueShort).Set(float64(qLen))
	return int(qLen)
}

// QueueWithNamespace returns the fully qualified name if Namespace is set, otherwise it just returns the queue name
func (rq *QueueMgr) QueueWithNamespace(queueShort string) string {
	if rq.clientOpts.Namespace != "" {
		return fmt.Sprintf("%s.%s", rq.clientOpts.Namespace, queueShort) // prepend namespace with dot
	} else {
		return queueShort
	}
}

// WaitForListenerShutdown blocking, waits until the wait group is down to 0
func (rq *QueueMgr) WaitForListenerShutdown() {
	rq.logger.Debug().Msgf("Wait for all listeners to exit. In a hurry? kill pid %d 😃", os.Getpid())
	rq.listenerWG.Wait()
	rq.logger.Debug().Msg("All queue listeners were shut down, goodbye")
}
