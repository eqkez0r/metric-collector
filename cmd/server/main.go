package main

import (
	"context"
	"github.com/Eqke/metric-collector/internal/config"
	httpserver "github.com/Eqke/metric-collector/internal/server"
	"github.com/Eqke/metric-collector/internal/storage/localstorage"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	sugLog := logger.Sugar()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	settings, err := config.NewServerConfig()
	if err != nil {
		log.Fatal(err)
	}
	storage := localstorage.New(sugLog)
	server := httpserver.New(ctx, settings, storage, sugLog)
	server.Run()
	server.Shutdown()
}
