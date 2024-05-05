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
	sugLog := logger.Sugar()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	settings, err := config.NewServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	//storage := localstorage.New(sugLog)
	//if settings.Restore {
	//	if err = storage.FromFile(settings.FileStoragePath); os.IsNotExist(err) {
	//		err = creatingStorageFile(settings, storage, sugLog)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//	}
	//
	//}
	//sugLog.Infof("Successful read from file")
	//
	//
	//if settings.DatabaseDSN != "" {
	//	conn, err = pgx.Connect(ctx, settings.DatabaseDSN)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	defer conn.Close(ctx)
	//}

	storage, err := storagemanager.GetStorage(ctx, sugLog, settings)
	if err != nil {
		log.Fatal(err)
	}

	server, err := httpserver.New(ctx, settings, storage, sugLog)
	if err != nil {
		log.Fatal(err)
	}
	server.Run()
	defer server.Shutdown()
}
