package config

import (
	"errors"
	"flag"
	e "github.com/Eqke/metric-collector/pkg/error"
	"github.com/ilyakaznacheev/cleanenv"
)

type AgentConfig struct {
	AgentEndpoint  string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

const (
	errPointNewAgentConfig = "error in NewAgentConfig(): "
)

func NewAgentConfig() (*AgentConfig, error) {
	cfg := &AgentConfig{}
	flag.StringVar(&cfg.AgentEndpoint, "a", defaultAddr, "agent endpoint")
	flag.IntVar(&cfg.ReportInterval, "r", defaultReportInterval, "report interval in seconds")
	flag.IntVar(&cfg.PollInterval, "p", defaultPollInterval, "poll interval in seconds")
	flag.Parse()
	if len(flag.Args()) != 0 {
		return nil, errors.New(errPointNewAgentConfig + errUnexpectedArguments.Error())
	}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, e.WrapError(errPointNewAgentConfig, err)
	}
	return cfg, nil
}
