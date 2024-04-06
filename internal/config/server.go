package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	Endpoint string
}

const (
	errPointNewServerConfig = "error in NewServerConfig(): "
)

var (
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
