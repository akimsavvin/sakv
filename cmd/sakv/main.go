package main

import (
	"context"
	"flag"
	"github.com/akimsavvin/sakv/internal/database/app"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "./config/config.yaml", "path to sakv config file")
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	if err := app.Start(ctx, configFilePath); err != nil {
		log.Fatalln(err)
	}
}
