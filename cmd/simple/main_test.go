package main

import (
	"os"
	"testing"
)

func TestMainFunc(t *testing.T) {

	os.Args = append(os.Args, "-name=tilltest")
	os.Args = append(os.Args, "-greeting=howdy")
	main() // will fail if name and greeting are unset

	// Test results here, and decide pass/fail.
}
