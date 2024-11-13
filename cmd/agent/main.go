package main

import (
	"context"
	"github.com/Eqke/metric-collector/internal/agent/config"
	"github.com/Eqke/metric-collector/internal/agent/grpcagent"
	"github.com/Eqke/metric-collector/internal/agent/httpagent"
	"github.com/Eqke/metric-collector/internal/agent/poller"
	"github.com/Eqke/metric-collector/internal/encrypting"
	"log"
	"os/signal"
	"sync"
	"syscall"

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
	sugarLogger := logger.Sugar()
	sugarLogger.Infoln(zap.String("Build version: ", buildVersion),
		zap.String("Build date: ", buildDate),
		zap.String("Git commit: ", buildCommit))
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	settings, err := config.NewAgentConfig()
	if err != nil {
		sugarLogger.Fatal(err)
	}

	publicKey, err := encrypting.GetPublicKey(settings.CryptoKey)
	if err != nil {
		sugarLogger.Fatal(err)
	}
	var wg sync.WaitGroup

	poll := poller.NewPoller(sugarLogger, settings)

	go poll.Poll(ctx, &wg)

	httpAgent := httpagent.New(settings, sugarLogger, publicKey, poll)

	wg.Add(2)
	go httpAgent.Run(ctx, &wg)

	grpcAgent := grpcagent.New(sugarLogger, settings, poll)

	go grpcAgent.Run(ctx, &wg)

	wg.Wait()
}
