package main

import (
	"context"
	"github.com/Eqke/metric-collector/internal/config"
	httpserver "github.com/Eqke/metric-collector/internal/server"
	"github.com/Eqke/metric-collector/internal/storage/localstorage"
	"go.uber.org/zap"
	"log"
	"os"
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
	if settings.Restore {
		if err = storage.FromFile(settings.FileStoragePath); os.IsNotExist(err) {
			err = creatingStorageFile(settings, storage, sugLog)
			if err != nil {
				log.Fatal(err)
			}
		}

	}
	sugLog.Infof("Successful read from file")
	server := httpserver.New(ctx, settings, storage, sugLog)
	server.Run()
	server.Shutdown()
}

func creatingStorageFile(
	settings *config.ServerConfig,
	storage *localstorage.LocalStorage,
	logger *zap.SugaredLogger) error {
	logger.Infof("File not found: %s", settings.FileStoragePath)
	logger.Info("Create new file")
	f, err := os.OpenFile(settings.FileStoragePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	b, err := storage.ToJSON()
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}
