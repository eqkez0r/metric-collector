package config

import "errors"

var (
	errUnexpectedArguments = errors.New("unexpected arguments")
)

const (
	defaultAddr           = "localhost:8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10

	EnvAddr           = "ADDRESS"
	EnvReportInterval = "REPORT_INTERVAL"
	EnvPollInterval   = "POLL_INTERVAL"
)
