// Package rpcclient see https://dev.to/karankumarshreds/go-rpc-implementation-4731
// and https://github.com/karankumarshreds/GoRPC/blob/master/client/client.go
package rpcclient

import (
	"log"
	"net/rpc"
)

type Args struct{}

func Run() int64 {
	// Address to this variable will be sent to the RPC server
	// Type of reply should be same as that specified on server
	var reply int64
	args := Args{}

	// DialHTTP connects to an HTTP RPC server at the specified network
	client, err := rpc.DialHTTP("tcp", "0.0.0.0:1234")
	if err != nil {
		log.Fatal("Client connection error: ", err)
	}

	// Invoke the remote function GiveServerTime attached to TimeServer pointer
	// Sending the arguments and reply variable address to the server as well
	err = client.Call("TimeServer.GiveServerTime", args, &reply)
	if err != nil {
		log.Fatal("Client invocation error: ", err)
	}

	// Print the reply from the server
	log.Printf("%d", reply)
	return reply
}
