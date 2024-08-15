// пакет config предоставляет структуру конфигурации для Agent
package config

import (
	"errors"
	"flag"
	e "github.com/Eqke/metric-collector/pkg/error"
	"github.com/ilyakaznacheev/cleanenv"
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
	AgentEndpoint  string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	HashKey        string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

// Функция NewAgentConfig создает экземлпяр типа AgentConfig
func NewAgentConfig() (*AgentConfig, error) {
	cfg := &AgentConfig{}
	flag.StringVar(&cfg.AgentEndpoint, "a", defaultAddr, "agent endpoint")
	flag.IntVar(&cfg.ReportInterval, "r", defaultReportInterval, "report interval in seconds")
	flag.IntVar(&cfg.PollInterval, "p", defaultPollInterval, "poll interval in seconds")
	flag.StringVar(&cfg.HashKey, "k", "", "hash key")
	flag.IntVar(&cfg.RateLimit, "l", defaultRateLimit, "rate limit")
	flag.Parse()
	if len(flag.Args()) != 0 {
		return nil, e.WrapError(errPointNewAgentConfig, ErrUnexpectedArguments)
	}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, e.WrapError(errPointNewAgentConfig, err)
	}
	return cfg, nil
}
