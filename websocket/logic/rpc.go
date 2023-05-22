package logic

import (
	"context"
	"flag"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"log"
	"time"
)

var Addr = flag.String("addr", "localhost:8972", "server address")

type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

func RpcClient() {
	flag.Parse()

	c := client.NewClient(client.DefaultOption)
	err := c.Connect("tcp", *Addr)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	args := &Args{
		A: 10,
		B: 20,
	}

	go func() {
		for {
			reply := &Reply{}
			err = c.Call(context.Background(), "Arith", "Mul", args, reply)

			if err != nil {
				log.Fatalf("failed to call: %v", err)
			}

			log.Printf("%d * %d = %d", args.A, args.B, reply.C)
			time.Sleep(3 * time.Second)
		}
	}()

	ch := make(chan *protocol.Message)
	c.RegisterServerMessageChan(ch)

	for msg := range ch {
		fmt.Printf("receive msg from server: %s\n", msg.Payload)
	}
}
