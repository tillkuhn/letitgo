/*
Copyright © 2022 Till Kuhn
*/
package main

import "github.com/tillkuhn/letitgo/cmd"

// main delegates to Cobra's Execute in cmd/root.go
func main() {
	cmd.Execute()
}
