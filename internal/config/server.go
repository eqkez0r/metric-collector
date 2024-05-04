package config

import (
	"flag"
	e "github.com/Eqke/metric-collector/pkg/error"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type ServerConfig struct {
	Endpoint        string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

const (
	errPointNewServerConfig = "error in NewServerConfig(): "
)

func NewServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}
	defaultStorePath := os.TempDir() + "/metrics-db.json"
	flag.StringVar(&cfg.Endpoint, "a", defaultAddr, "server endpoint")
	flag.IntVar(&cfg.StoreInterval, "i", defaultStoreInterval, "store interval in seconds")
	flag.StringVar(&cfg.FileStoragePath, "f", defaultStorePath, "file storage path")
	flag.BoolVar(&cfg.Restore, "r", defaultRestoreVal, "restore")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "database DSN")
	flag.Parse()

	if len(flag.Args()) != 0 {
		return nil, errUnexpectedArguments
	}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, e.WrapError(errPointNewServerConfig, err)
	}
	return cfg, nil
}
