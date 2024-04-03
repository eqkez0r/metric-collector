package main

import (
	"context"
	"flag"
	"github.com/Eqke/metric-collector/internal/agent"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	flagAddr           string
	flagReportInterval int
	flagPollInterval   int
)

const (
	defaultAddr           = "localhost:8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10
)

type EnvCfg struct {
	address        string `env:"ADDRESS"`
	reportInterval int    `env:"REPORT_INTERVAL"`
	pollInterval   int    `env:"POLL_INTERVAL"`
}

func parseFlags() {
	flag.StringVar(&flagAddr, "a", defaultAddr, "agent endpoint")
	flag.IntVar(&flagReportInterval, "r", defaultReportInterval, "report interval in seconds")
	flag.IntVar(&flagPollInterval, "p", defaultPollInterval, "poll interval in seconds")
	flag.Parse()
	if len(flag.Args()) != 0 {
		log.Println("unexpected arguments: ", flag.Args())
		os.Exit(1)
	}
	if v, ok := os.LookupEnv("ADDRESS"); ok {
		flagAddr = v
	}
	if v, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		flagReportInterval, _ = strconv.Atoi(v)
	}
	if v, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		flagPollInterval, _ = strconv.Atoi(v)
	}
	//var cfg EnvCfg
	//err := env.Parse(&cfg)
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Println(cfg)
	//if cfg.address != "" {
	//	flagAddr = cfg.address
	//}
	//if cfg.reportInterval != 0 {
	//	flagReportInterval = cfg.reportInterval
	//}
	//if cfg.pollInterval != 0 {
	//	flagPollInterval = cfg.pollInterval
	//}
}

func main() {
	parseFlags()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	a := agent.New(ctx, flagAddr, flagPollInterval, flagReportInterval)
	a.Run()
}
