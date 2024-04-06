package config

import (
	"errors"
	"flag"
	"os"
)

type ServerConfig struct {
	Endpoint string
}

var (
	errUnexpectedArguments = errors.New("unexpected arguments")

	flagServerAddr string
)

func NewServerConfig() (*ServerConfig, error) {
	flag.StringVar(&flagServerAddr, "a", defaultAddr, "server endpoint")
	flag.Parse()

	if len(flag.Args()) != 0 {
		return nil, errUnexpectedArguments
	}
	if v, ok := os.LookupEnv(EnvAddr); ok {
		flagServerAddr = v
	}

	return &ServerConfig{
		Endpoint: flagServerAddr,
	}, nil
}
