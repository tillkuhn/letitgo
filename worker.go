package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jpillora/backoff"
)

// Worker is a unit that performs the task
type Worker struct {
	name      string // name is a Human friendly worker name
	processed int    // processed count to track  messages per worker
}

const maxAttempts = 3

// Exponential Backoff based on https://github.com/jpillora/backoff
func main() {
	w := Worker{
		name: "hase",
	}
	log.Println(w.Status())
	w.HandleJob()
	w.HandleJob()
	w.HandleJob()
	log.Println(w.Status())
	// runWithRetry()

	// Add resp. upsert document(s)
	document := struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Text string `json:"text"`
	}{
		ID:   "1000",
		Name: "Go",
		Text: "Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.",
	}
	log.Printf("%v", document)
	runWithRetry()
}

// HandleJob increments the processed counter
// Methods can be defined for either pointer or value receiver, go handles conversion automatically
// We use a pointer receiver type to avoid copying on method calls
// or to allow the method to mutate the receiving struct.
func (w *Worker) HandleJob() {
	log.Printf("Working %d\n", w.processed)
	w.processed = w.processed + 1
}

func (w Worker) Status() string {
	return fmt.Sprintf("%s: count %d", w.name, w.processed)
}

func runWithRetry() {
	attempts := 0
	b := &backoff.Backoff{
		//These are the defaults
		Min:    1 * time.Second,
		Max:    10 * time.Second,
		Factor: 2,
		Jitter: false,
	}
	var conn net.Conn
	for {
		attempts++
		log.Printf("Attempt #%d/%d: Trying ....\n", attempts, maxAttempts)
		var err error
		conn, err = net.DialTimeout("tcp", "example.com:5309", 2*time.Second)
		if err == nil {
			log.Printf("Connected after %d attempt(s)", attempts)
			break
		} else if attempts >= maxAttempts {
			log.Fatalf("%s, max Attemps %d reached. I give up", err, maxAttempts)
		} else if attempts > 1 {
			d := b.Duration()
			log.Printf("%s, reconnecting in %s\n", err, d)
			time.Sleep(d)
		}
	}
	//connected
	b.Reset()
	if _, err := conn.Write([]byte("hello world!")); err != nil {
		log.Println(err)
	}
	// ... Read ... Write ... etc
	conn.Close()
	//disconnected
}
