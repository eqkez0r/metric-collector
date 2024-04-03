package main

import (
	"flag"
	"log"
	"os"
)

var (
	flagAddr           string
	flagReportInterval int
	flagPollInterval   int
)

func parseFlags() {
	flag.StringVar(&flagAddr, "a", "localhost:8080", "agent endpoint")
	flag.IntVar(&flagReportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&flagPollInterval, "p", 2, "poll interval in seconds")
	flag.Parse()
	if len(flag.Args()) != 0 {
		log.Println("unexpected arguments: ", flag.Args())
		os.Exit(1)
	}
}
