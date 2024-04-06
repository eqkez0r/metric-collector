package config

import (
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
)

type AgentConfig struct {
	AgentEndpoint  string
	ReportInterval int
	PollInterval   int
}

const (
	errPointNewAgentConfig = "error in NewAgentConfig(): "
)

var (
	flagAgentAddr      string
	flagReportInterval int
	flagPollInterval   int
)

func NewAgentConfig() (*AgentConfig, error) {
	flag.StringVar(&flagAgentAddr, "a", defaultAddr, "agent endpoint")
	flag.IntVar(&flagReportInterval, "r", defaultReportInterval, "report interval in seconds")
	flag.IntVar(&flagPollInterval, "p", defaultPollInterval, "poll interval in seconds")
	flag.Parse()
	if len(flag.Args()) != 0 {
		return nil, errors.New(errPointNewAgentConfig + errUnexpectedArguments.Error())
	}
	if v, ok := os.LookupEnv(EnvAddr); ok {
		log.Println("ADDRESS", v)
		flagAgentAddr = v
	}
	if v, ok := os.LookupEnv(EnvReportInterval); ok {
		if reportInterval, err := strconv.Atoi(v); err != nil {
			return nil, errors.New(errPointNewAgentConfig + errUnexpectedArguments.Error())
		} else {
			flagReportInterval = reportInterval
		}
	}
	if v, ok := os.LookupEnv(EnvPollInterval); ok {
		if pollInterval, err := strconv.Atoi(v); err != nil {
			return nil, err
		} else {
			flagPollInterval = pollInterval
		}
	}
	return &AgentConfig{
		AgentEndpoint:  flagAgentAddr,
		ReportInterval: flagReportInterval,
		PollInterval:   flagPollInterval,
	}, nil
}
