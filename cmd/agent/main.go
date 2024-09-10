package main

import (
	"context"
	"github.com/Eqke/metric-collector/internal/agent/config"
	"github.com/Eqke/metric-collector/internal/encrypting"
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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	settings, err := config.NewAgentConfig()
	if err != nil {
		sugLog.Fatal(err)
	}

	publicKey, err := encrypting.GetPublicKey(settings.CryptoKey)
	if err != nil {
		sugLog.Fatal(err)
	}

	a := agent.New(settings, sugLog, publicKey)
	a.Run(ctx)
}
