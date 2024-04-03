package main

import (
	"flag"
	"log"
	"os"
)

var flagAddr string

func parseFlags() {
	flag.StringVar(&flagAddr, "a", "localhost:8080", "server endpoint")

	flag.Parse()

	if len(flag.Args()) != 0 {
		log.Println("unexpected arguments: ", flag.Args())
		os.Exit(1)
	}
}
