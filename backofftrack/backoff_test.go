package backofftrack

import (
	"fmt"
	"testing"
)

func TestReconnect(t *testing.T) {
	// if launched with 'go test -short   ./...'
	if testing.Short() {
		fmt.Printf("Short test, using connectWithRetry with very low values\n")
		maxAttempts = 1
	}
	ConnectWithRetry()
}
