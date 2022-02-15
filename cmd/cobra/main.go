package main

import (
	"log"

	"github.com/spf13/cobra/cobra/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
