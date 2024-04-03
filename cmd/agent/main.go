package main

import (
	"github.com/Eqke/metric-collector/internal/agent"
	"time"
)

func main() {
	a := agent.New("localhost:8080", 5*time.Second)
	a.Run()
}
