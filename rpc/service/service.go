package service

import (
	"context"
)

var (
	Addr     = "localhost:8972"
	EtcdAddr = "localhost:2379"
	BasePath = "/rpcx_test"
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
