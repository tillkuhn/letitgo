package main

import (
	"github.com/spf13/cobra/cobra/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
