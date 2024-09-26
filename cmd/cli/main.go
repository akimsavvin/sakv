package main

import (
	"context"
	"flag"
	"net"
	"os"
	"os/signal"
	"sakv/internal/cli/listener"
	"sakv/internal/cli/query"
	"syscall"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", "localhost:3223", "sakv database server address")
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}

	qh := query.NewNetHandler(conn)
	l := listener.New(qh)

	l.StartListening(ctx)
}
