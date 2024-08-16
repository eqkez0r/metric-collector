// Пакет storagemanager предназначен для описания фабрики хранилищ.
package storagemanager

import (
	"context"
	"github.com/Eqke/metric-collector/internal/server/config"
	"github.com/Eqke/metric-collector/internal/storage"
	"os"

	"github.com/Eqke/metric-collector/internal/storage/localstorage"
	"github.com/Eqke/metric-collector/internal/storage/postgres"
	"go.uber.org/zap"
)

// Объявление ошибок
const (
	ErrPointGetStorage = "error in storagemanager.GetStorage(): "
)

// Функция GetStorage, которая получает ctx - контекст,
// logger - для логирования,
// cfg - для получения данных из конфигурации приложения.
func GetStorage(ctx context.Context, logger *zap.SugaredLogger, cfg *config.ServerConfig) (storage.Storage, error) {
	switch {
	case cfg.DatabaseDSN != "":
		{
			return postgres.New(ctx, logger, cfg.DatabaseDSN)
		}
	default:
		{
			s := localstorage.New(logger)
			if cfg.Restore {
				if err := s.FromFile(ctx, cfg.FileStoragePath); os.IsNotExist(err) {
					err = creatingStorageFile(ctx, cfg, s, logger)
					if err != nil {
						logger.Fatalf("%v: %v", ErrPointGetStorage, err)
					}
				} else {
					err = s.FromFile(ctx, cfg.FileStoragePath)
					if err != nil {
						logger.Fatalf("%v: %v", ErrPointGetStorage, err)
					}
				}
			}
			logger.Info("Successful read from file")
			return s, nil
		}
	}
}

// Функция creatingStorageFile используется для backup в memory
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
