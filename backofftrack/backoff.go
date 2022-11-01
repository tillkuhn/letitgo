package backofftrack

import (
	"fmt"
	"net"
	"time"

	"github.com/jpillora/backoff"
)

var (
	target      = "example.com:530933"
	maxAttempts = 4.0
)

// ConnectWithRetry implements https://github.com/jpillora/backoff
func ConnectWithRetry() {
	b := &backoff.Backoff{
		Min:    1 * time.Millisecond,
		Max:    1 * time.Second,
		Factor: 2,
		Jitter: false,
	}

	var conn net.Conn
	var err error
	for {
		fmt.Printf("Attempt %f, connecting in %s\n", b.Attempt(), target)
		d := net.Dialer{Timeout: 500 * time.Millisecond}
		conn, err = d.Dial("tcp", target)
		if err != nil {
			// Attempt returns the current attempt counter value.
			if b.Attempt() >= maxAttempts-1 {
				fmt.Printf("max attempts %f reached, give up\n", maxAttempts)
				return
			}
			// Duration returns the duration for the current attempt before incrementing
			sleepy := b.Duration()

			fmt.Printf("Error %s, reconnecting in %v\n", err, sleepy)
			time.Sleep(sleepy)
		} else {
			break
		}
	}
	//connected
	b.Reset()
	if conn != nil {
		_, _ = conn.Write([]byte("hello world!"))
		// ... Read ... Write ... etc
		_ = conn.Close()
	}

}
