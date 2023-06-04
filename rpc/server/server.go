package main

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx/server"
	"im/rpc/service"
	"net"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
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
	quit := make(chan os.Signal, 1) // 退出信号
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	srv := server.NewServer()
	//s.RegisterName("Arith", new(service.Arith), "")
	srv.Register(new(Arith), "")
	go srv.Serve("tcp", service.Addr)
	go shutdown(quit, srv)

	for !connected {
		time.Sleep(time.Second)
	}

	fmt.Printf("start to send messages to %s\n", clientConn.RemoteAddr().String())
	for {
		if clientConn != nil {
			err := srv.SendMessage(clientConn, "", "", nil, []byte("ping"))
			if err != nil {
				fmt.Printf("failed to send messsage to %s: %v\n", clientConn.RemoteAddr().String(), err)
				clientConn = nil
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func shutdown(quit chan os.Signal, srv *server.Server) {
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("server shutdown...")
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
	os.Exit(0)
}
