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
	parseFlags()

	settings := &httpserver.Settings{
		Endpoint: flagAddr,
	}
	server := httpserver.New(ctx, settings)

	server.Run()

}
