package main

import (
	"context"
	"github.com/Eqke/metric-collector/internal/config"
	httpserver "github.com/Eqke/metric-collector/internal/server"
	"github.com/Eqke/metric-collector/internal/storagemanager"
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
	sugarLogger := logger.Sugar()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	settings, err := config.NewServerConfig()
	if err != nil {
		sugarLogger.Fatal(err)
	}

	storage, err := storagemanager.GetStorage(ctx, sugarLogger, settings)
	if err != nil {
		sugarLogger.Fatal(err)
	}

	server, err := httpserver.New(ctx, settings, storage, sugarLogger)
	if err != nil {
		sugarLogger.Fatal(err)
	}
	server.Run(ctx)

}
