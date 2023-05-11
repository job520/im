package service

import (
	"context"
	"flag"
)

var (
	Addr     = flag.String("addr", "localhost:8972", "server address")
	EtcdAddr = flag.String("etcdAddr", "localhost:2379", "etcd address")
	BasePath = flag.String("base", "/rpcx_test", "prefix path")
)

type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

type Arith int

func (t *Arith) Mul(ctx context.Context, args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}
