package main

import (
	"context"
	"github.com/Eqke/metric-collector/internal/encrypting"
	"github.com/Eqke/metric-collector/internal/restorer"
	"github.com/Eqke/metric-collector/internal/server/config"
	"github.com/Eqke/metric-collector/internal/server/grpcserver"
	"github.com/Eqke/metric-collector/internal/server/httpserver"
	"log"
	"os/signal"
	"sync"
	"syscall"

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
	var wg sync.WaitGroup
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	sugarLogger := logger.Sugar()
	sugarLogger.Infoln(
		zap.String("Build version: ", buildVersion),
		zap.String("Build date: ", buildDate),
		zap.String("Git commit: ", buildCommit))
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
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

	if settings.Restore && settings.DatabaseDSN == "" {
		restore := restorer.New(sugarLogger, storage, settings.FileStoragePath, settings.StoreInterval)
		wg.Add(1)
		go restore.Run(ctx, &wg)
	}

	server := httpserver.New(settings, storage, sugarLogger, privateKey)
	wg.Add(2)
	go server.Run(ctx, &wg)

	grpcServer := grpcserver.New(sugarLogger, storage, settings.GrpcServerHost)

	go grpcServer.Run(ctx, &wg)

	wg.Wait()
}
