// Пакет config содержит конфигурацию сервера
package config

import (
	"errors"
	"flag"
	e "github.com/Eqke/metric-collector/pkg/error"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

const (
	// Объявление точки ошибки
	errPointNewServerConfig = "error in NewServerConfig(): "

	// Значение хоста по умолчанию
	defaultAddr = "localhost:8080"
	// Значение периода записи в файл по умолчанию
	defaultStoreInterval = 300
	// Значение флага, отвечающий за запись данных в файл
	defaultRestoreVal = true
)

var (
	// Объявление ошибки об недопустимых переменных
	ErrUnexpectedArguments = errors.New("unexpected arguments")
)

// Тип ServerConfig представляет структуру для конфигурации сервера
// Host - адрес и порт, на котором крутится сервер
type ServerConfig struct {
	Host            string `env:"ADDRESS" json:"address"`
	StoreInterval   int    `env:"STORE_INTERVAL" json:"store_interval"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"store_file"`
	Restore         bool   `env:"RESTORE" json:"restore"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	HashKey         string `env:"KEY"`
	CryptoKey       string `env:"CRYPTO_KEY" json:"crypto_key"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
}

// Функция NewServerConfig возвращает экземпляр конфигурации сервера
func NewServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}
	var cfgPathFl string
	defaultStorePath := os.TempDir() + "/metrics-db.json"
	flag.StringVar(&cfg.Host, "a", defaultAddr, "server host")
	flag.IntVar(&cfg.StoreInterval, "i", defaultStoreInterval, "store interval in seconds")
	flag.StringVar(&cfg.FileStoragePath, "f", defaultStorePath, "file storage path")
	flag.BoolVar(&cfg.Restore, "r", defaultRestoreVal, "restore")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "Database DSN")
	flag.StringVar(&cfg.HashKey, "k", "", "hash key")
	flag.StringVar(&cfg.CryptoKey, "s", "", "path to crypto key")
	flag.StringVar(&cfg.TrustedSubnet, "t", "", "trusted subnet (CIDR)")
	flag.StringVar(&cfgPathFl, "c", "", "path to cfg")

	flag.Parse()

	if configPath := os.Getenv("CONFIG"); configPath != "" {
		cfgPathFl = configPath
	}

	if cfgPathFl != "" {
		err := cleanenv.ReadConfig(cfgPathFl, &cfg)
		if err != nil {
			return nil, e.WrapError(errPointNewServerConfig, err)
		}
	}
	if len(flag.Args()) != 0 {
		return nil, e.WrapError(errPointNewServerConfig, ErrUnexpectedArguments)
	}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, e.WrapError(errPointNewServerConfig, err)
	}

	return cfg, nil
}
