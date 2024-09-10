package main

import (
	"context"
	"github.com/Eqke/metric-collector/internal/encrypting"
	"github.com/Eqke/metric-collector/internal/server/config"
	"log"
	"os/signal"
	"syscall"

	httpserver "github.com/Eqke/metric-collector/internal/server"
	"github.com/Eqke/metric-collector/internal/storagemanager"
	"go.uber.org/zap"

	_ "net/http/pprof"
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
	sugarLogger.Infoln(
		zap.String("Build version: ", buildVersion),
		zap.String("Build date: ", buildDate),
		zap.String("Git commit: ", buildCommit))
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

	privateKey, err := encrypting.GetPrivateKey(settings.CryptoKey)
	if err != nil {
		sugarLogger.Error(err)
		err = encrypting.GenerateIfNotExist(settings.CryptoKey)
		if err != nil {
			sugarLogger.Fatal(err)
		}

		privateKey, err = encrypting.GetPrivateKey(settings.CryptoKey)
		if err != nil {
			sugarLogger.Fatal(err)
		}
	}

	server, err := httpserver.New(ctx, settings, storage, sugarLogger, privateKey)
	if err != nil {
		sugarLogger.Fatal(err)
	}
	server.Run(ctx)

}
