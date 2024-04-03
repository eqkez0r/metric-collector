package main

import (
	"context"
	httpserver "github.com/Eqke/metric-collector/internal/server"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	settings := &httpserver.Settings{
		Host: "0.0.0.0",
		Port: "8080",
	}
	server := httpserver.New(ctx, settings)

	server.Run()

}
