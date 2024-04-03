package main

import (
	"context"
	"flag"
	httpserver "github.com/Eqke/metric-collector/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var flagAddr string

const (
	defaultAddr = "localhost:8080"
)

type EnvCfg struct {
	address string `env:"ADDRESS"`
}

func parseFlags() {
	flag.StringVar(&flagAddr, "a", defaultAddr, "server endpoint")

	flag.Parse()

	if len(flag.Args()) != 0 {
		log.Println("unexpected arguments: ", flag.Args())
		os.Exit(1)
	}
	if v, ok := os.LookupEnv("ADDRESS"); ok {
		flagAddr = v
	}

}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	parseFlags()

	settings := &httpserver.Settings{
		Endpoint: flagAddr,
	}
	server := httpserver.New(ctx, settings)

	server.Run()

}
