// Package rpc see https://dev.to/karankumarshreds/go-rpc-implementation-4731
// and https://github.com/karankumarshreds/GoRPC
package rpc

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	sc      = make(chan os.Signal, 1)
	running = make(chan bool, 1)
)

type Args struct{}

type TimeServer int64

func Serve() {
	timeserver := new(TimeServer)
	// Register the timeserver object upon which the GiveServerTime
	// function will be called from the RPC server (from the client)
	if err := rpc.Register(timeserver); err != nil {
		log.Fatal(err)
	}
	// Registers an HTTP handler for RPC messages
	rpc.HandleHTTP() // ?
	// Start listening for the requests on port 1234
	listener, err := net.Listen("tcp", "0.0.0.0:1234")
	if err != nil {
		log.Fatal("Listener error: ", err)
	}
	// Serve accepts incoming HTTP connections on the listener l, creating
	// a new service goroutine for each. The service goroutines read requests
	// and then call handler to reply to them
	// https://stackoverflow.com/questions/61571993/how-to-shutdown-a-rpc-server-in-golang
	log.Printf("Starting listener %v \n", listener.Addr())
	server := http.Server{}
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()
	running <- true
	log.Println("Waiting for customers")
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
	_ = server.Close()
	log.Println("Goodbye from rpc")
}

func (t *TimeServer) GiveServerTime(_ *Args, reply *int64) error {
	// Set the value at the pointer got from the client
	log.Println("Time requested")
	*reply = time.Now().Unix()
	return nil
}
