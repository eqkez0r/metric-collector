// Пакет config содержит конфигурацию сервера
package config

import (
	"errors"
	"flag"
	"os"

	e "github.com/Eqke/metric-collector/pkg/error"
	"github.com/ilyakaznacheev/cleanenv"
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
	Host            string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	HashKey         string `env:"KEY"`
}

// Функция NewServerConfig возвращает экземпляр конфигурации сервера
func NewServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}
	defaultStorePath := os.TempDir() + "/metrics-db.json"
	flag.StringVar(&cfg.Host, "a", defaultAddr, "server host")
	flag.IntVar(&cfg.StoreInterval, "i", defaultStoreInterval, "store interval in seconds")
	flag.StringVar(&cfg.FileStoragePath, "f", defaultStorePath, "file storage path")
	flag.BoolVar(&cfg.Restore, "r", defaultRestoreVal, "restore")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "Database DSN")
	flag.StringVar(&cfg.HashKey, "k", "", "hash key")
	flag.Parse()

	if len(flag.Args()) != 0 {
		return nil, e.WrapError(errPointNewServerConfig, ErrUnexpectedArguments)
	}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, e.WrapError(errPointNewServerConfig, err)
	}

	return cfg, nil
}
