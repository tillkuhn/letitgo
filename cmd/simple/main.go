package main

import (
	"flag"
	"log"
	"os"
)

// Call DEBUG=true go run cmd/mycli/main.go  -name hase  -prompt

func main() {
	var name string
	var greeting string
	var prompt bool
	var preview bool

	flag.StringVar(&name,"name","","Your name")
	flag.StringVar(&greeting,"greeting","","Greetings")
	flag.BoolVar(&prompt,"prompt",false,"Shall we prompt?")
	flag.BoolVar(&preview,"preview",false,"Use preview")

	flag.Parse()

	if ! prompt  && (name == "" || greeting == "") {
		flag.Usage()
		os.Exit(1)
	}

	if os.Getenv("DEBUG") != "" {
		log.Printf("[DEBUG] Name: %s",name)
	}
	log.Printf("Welcome to my CLI, %s",name)
}
