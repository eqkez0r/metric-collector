package main

import (
	"context"
	"github.com/Eqke/metric-collector/internal/agent"
	"os/signal"
	"syscall"
)

func main() {
	parseFlags()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	a := agent.New(ctx, flagAddr, flagPollInterval, flagReportInterval)
	a.Run()
}
