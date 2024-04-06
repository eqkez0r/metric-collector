package config

import (
	"flag"
	e "github.com/Eqke/metric-collector/pkg/error"
	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Endpoint string `env:"ADDRESS"`
}

const (
	errPointNewServerConfig = "error in NewServerConfig(): "
)

func NewServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}
	flag.StringVar(&cfg.Endpoint, "a", defaultAddr, "server endpoint")
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
