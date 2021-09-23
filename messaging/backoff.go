package messaging

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jpillora/backoff"
)

func RunWithBackoff(maxAttempts int) error {
	attempts := 0
	b := &backoff.Backoff{
		//These are the defaults
		Min:    10 * time.Millisecond,
		Max:    2000 * time.Millisecond,
		Factor: 2,
		Jitter: false,
	}
	var conn net.Conn
	for {
		attempts++
		log.Printf("Attempt #%d/%d: Trying ....\n", attempts, maxAttempts)
		var err error
		conn, err = net.DialTimeout("tcp", "no.such.host:5309", 1*time.Second)
		if err == nil {
			log.Printf("Connected after %d attempt(s)", attempts)
			break
			// we have an error - either try again or give up
		} else if attempts >= maxAttempts {
			errMsg := fmt.Sprintf("error %s, max Attemps %d reached. I give up", err, maxAttempts)
			log.Print(errMsg)
			return errors.New(errMsg)
		} else {
			d := b.Duration()
			log.Printf("error %s, reconnecting in %s", err, d)
			time.Sleep(d)
		}
	}
	//connected
	b.Reset()
	// do something with con
	err := conn.Close()
	return err
}
