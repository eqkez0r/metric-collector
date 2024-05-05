package storagemanager

import (
	"context"
	"github.com/Eqke/metric-collector/internal/config"
	stor "github.com/Eqke/metric-collector/internal/storage"
	"github.com/Eqke/metric-collector/internal/storage/localstorage"
	"github.com/Eqke/metric-collector/internal/storage/postgres"
	"go.uber.org/zap"
)

func GetStorage(ctx context.Context, logger *zap.SugaredLogger, cfg *config.ServerConfig) (stor.Storage, error) {
	if cfg.DatabaseDSN != "" {
		return postgres.New(ctx, logger, cfg.DatabaseDSN)
	}
	return localstorage.New(logger), nil
}
