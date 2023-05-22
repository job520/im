package main

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx/server"
	"im/rpc/service"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var clientConn net.Conn
var connected = false

type Arith int

func (t *Arith) Mul(ctx context.Context, args *service.Args, reply *service.Reply) error {
	clientConn = ctx.Value(server.RemoteConnContextKey).(net.Conn)
	reply.C = args.A * args.B
	connected = true
	return nil
}

func main() {

	ln, _ := net.Listen("tcp", ":9981")
	go http.Serve(ln, nil)

	s := server.NewServer()
	//s.RegisterName("Arith", new(service.Arith), "")
	s.Register(new(Arith), "")
	go s.Serve("tcp", service.Addr)

	for !connected {
		time.Sleep(time.Second)
	}

	fmt.Printf("start to send messages to %s\n", clientConn.RemoteAddr().String())
	for {
		if clientConn != nil {
			err := s.SendMessage(clientConn, "", "", nil, []byte("ping"))
			if err != nil {
				fmt.Printf("failed to send messsage to %s: %v\n", clientConn.RemoteAddr().String(), err)
				clientConn = nil
			}
		}
		time.Sleep(10 * time.Second)
	}
}
