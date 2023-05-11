package main

import (
	"context"
	"flag"
	"fmt"
	"im/rpc/service"
	"log"

	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

func main() {
	flag.Parse()

	c := client.NewClient(client.DefaultOption)
	err := c.Connect("tcp", *service.Addr)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	args := &service.Args{
		A: 10,
		B: 20,
	}

	// invoke once and set up the connection
	reply := &service.Reply{}
	err = c.Call(context.Background(), "Arith", "Mul", args, reply)

	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Printf("%d * %d = %d", args.A, args.B, reply.C)

	ch := make(chan *protocol.Message)
	c.RegisterServerMessageChan(ch)

	for msg := range ch {
		fmt.Printf("receive msg from server: %s\n", msg.Payload)
	}
}
