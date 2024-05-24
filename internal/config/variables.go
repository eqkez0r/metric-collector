package config

import (
	"errors"
	"os"
)

var (
	errUnexpectedArguments = errors.New("unexpected arguments")
	tmpdir                 = os.TempDir()
)

const (
	defaultAddr           = "localhost:8080"
	defaultStoreInterval  = 300
	defaultRestoreVal     = true
	defaultPollInterval   = 2
	defaultReportInterval = 10
	defaultRateLimit      = 100
)
