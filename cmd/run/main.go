package main

import (
	"log"
	"os"
	"tillkuhn/goplay/concurrent"
	"tillkuhn/goplay/monitor"
)

func main() {
	cmd := os.Args[1]
	if cmd == "job-queue" {
		concurrent.DoWork()
	} else if cmd == "prometheus" {
		monitor.RunPrometheus()
	} else {
		log.Fatal("Unknown or no cmd " + cmd)
	}
}
