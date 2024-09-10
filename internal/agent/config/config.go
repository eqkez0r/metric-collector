// пакет config предоставляет структуру конфигурации для Agent
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
	errPointNewAgentConfig = "error in NewAgentConfig(): "

	// Значение адреса сервера по умолчанию
	defaultAddr = "localhost:8080"
	// Значение таймера для обновления метрик
	defaultPollInterval = 2
	// Значение таймера для публикации метрик
	defaultReportInterval = 10
	// Ограничение на кол-во запросов по умолчанию
	defaultRateLimit = 100
)

var (
	ErrUnexpectedArguments = errors.New("unexpected arguments")
)

// Тип AgentConfig является типом конфигурации для Agent
type AgentConfig struct {
	AgentEndpoint  string `env:"ADDRESS" json:"address"`
	ReportInterval int    `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval   int    `env:"POLL_INTERVAL" json:"poll_interval"`
	HashKey        string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
	CryptoKey      string `env:"CRYPTO_KEY" json:"crypto_key"`
}

// Функция NewAgentConfig создает экземлпяр типа AgentConfig
func NewAgentConfig() (*AgentConfig, error) {
	cfg := &AgentConfig{}
	var cfgPathFl string

	flag.StringVar(&cfg.AgentEndpoint, "a", defaultAddr, "agent endpoint")
	flag.IntVar(&cfg.ReportInterval, "r", defaultReportInterval, "report interval in seconds")
	flag.IntVar(&cfg.PollInterval, "p", defaultPollInterval, "poll interval in seconds")
	flag.StringVar(&cfg.HashKey, "k", "", "hash key")
	flag.IntVar(&cfg.RateLimit, "l", defaultRateLimit, "rate limit")
	flag.StringVar(&cfg.CryptoKey, "s", "", "path to crypto key")
	flag.StringVar(&cfgPathFl, "c", "", "path to cfg")
	flag.Parse()

	if configPath := os.Getenv("CONFIG"); configPath != "" {
		cfgPathFl = configPath
	}

	if cfgPathFl != "" {
		err := cleanenv.ReadConfig(cfgPathFl, &cfg)
		if err != nil {
			return nil, e.WrapError(errPointNewAgentConfig, err)
		}
	}
	if len(flag.Args()) != 0 {
		return nil, e.WrapError(errPointNewAgentConfig, ErrUnexpectedArguments)
	}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, e.WrapError(errPointNewAgentConfig, err)
	}

	return cfg, nil
}
