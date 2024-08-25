package main

import (
	"context"
	"github.com/Eqke/metric-collector/internal/agent/config"
	"log"
	"os/signal"
	"syscall"

	"github.com/Eqke/metric-collector/internal/agent"
	"go.uber.org/zap"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	sugLog := logger.Sugar()
	sugLog.Infoln(zap.String("Build version: ", buildVersion),
		zap.String("Build date: ", buildDate),
		zap.String("Git commit: ", buildCommit))
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	settings, err := config.NewAgentConfig()
	if err != nil {
		log.Fatal(err)
	}
	a := agent.New(settings, sugLog)
	a.Run(ctx)
}
