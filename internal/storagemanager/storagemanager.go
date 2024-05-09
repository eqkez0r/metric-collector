package storagemanager

import (
	"context"
	"github.com/Eqke/metric-collector/internal/config"
	stor "github.com/Eqke/metric-collector/internal/storage"
	"github.com/Eqke/metric-collector/internal/storage/localstorage"
	"github.com/Eqke/metric-collector/internal/storage/postgres"
	"go.uber.org/zap"
	"os"
)

const (
	ErrPointGetStorage = "error in storagemanager.GetStorage(): "
)

func GetStorage(ctx context.Context, logger *zap.SugaredLogger, cfg *config.ServerConfig) (stor.Storage, error) {
	if cfg.DatabaseDSN != "" {
		return postgres.New(ctx, logger, cfg.DatabaseDSN)
	}
	storage := localstorage.New(logger)
	if cfg.Restore {
		if err := storage.FromFile(ctx, cfg.FileStoragePath); os.IsNotExist(err) {
			err = creatingStorageFile(ctx, cfg, storage, logger)
			if err != nil {
				logger.Fatalf("%v: %v", ErrPointGetStorage, err)
			}
		} else {
			err = storage.FromFile(ctx, cfg.FileStoragePath)
			if err != nil {
				logger.Fatalf("%v: %v", ErrPointGetStorage, err)
			}
		}
	}
	logger.Infof("Successful read from file")
	return storage, nil
}

func creatingStorageFile(
	ctx context.Context,
	settings *config.ServerConfig,
	storage *localstorage.LocalStorage,
	logger *zap.SugaredLogger) error {
	logger.Infof("File not found: %s", settings.FileStoragePath)
	logger.Info("Create new file")
	f, err := os.OpenFile(settings.FileStoragePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	b, err := storage.ToJSON(ctx)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}
