package agent

import (
	"github.com/Eqke/metric-collector/pkg/metric"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	pollCounterName = "pollCounter"
	randName        = "random"
	typeGauge       = "gauge"
	typeCounter     = "counter"
)

type Agent struct {
	addr         string
	client       http.Client
	pollInterval time.Duration
	pollCounter  int64
}

func New(addr string, pollInterval time.Duration) *Agent {
	return &Agent{
		addr:         addr,
		client:       http.Client{},
		pollInterval: pollInterval,
		pollCounter:  0,
	}
}

func (a *Agent) Run() {
	log.Println("agent was started")
	ms := &runtime.MemStats{}
	m := metric.PrepareMetrics(ms)

	for range time.Tick(a.pollInterval) {
		a.pollCounter++
		m[typeCounter][pollCounterName] = strconv.FormatInt(a.pollCounter, 10)
		metric.UpdateMetrics(ms, m)
		for metricType, metricMap := range m {
			for metricName, metricValue := range metricMap {
				url := strings.Join([]string{"http:/", a.addr, "update", metricType.String(), metricName.String(), metricValue}, "/")
				resp, err := a.client.Post(url, "text/plain", nil)
				if err != nil {
					log.Println(err)
				}
				resp.Body.Close()
			}
		}
	}
}

func (a *Agent) Stop() {

}

func (a *Agent) GetPollCounter() int64 {
	return a.pollCounter
}
