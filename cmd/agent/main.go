package main

import (
	"context"
	"github.com/Eqke/metric-collector/internal/agent"
	"github.com/Eqke/metric-collector/internal/config"
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
	settings, err := config.NewAgentConfig()
	if err != nil {
		log.Fatal(err)
	}
	a := agent.New(ctx, settings, sugLog)
	a.Run()
}
