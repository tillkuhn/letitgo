package httpserver

import "testing"

func TestRateLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test")
	}
	Run()
}
