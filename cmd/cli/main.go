package main

import (
	"context"
	"flag"
	"github.com/akimsavvin/sakv/internal/cli/listener"
	"github.com/akimsavvin/sakv/internal/cli/query"
	"net"
	"os"
	"os/signal"
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
