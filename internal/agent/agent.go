package agent

import (
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/go-resty/resty/v2"
	"log"
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
	client       *resty.Client
	pollInterval time.Duration
	pollCounter  int64
}

func New(addr string, pollInterval time.Duration) *Agent {
	return &Agent{
		addr:         addr,
		client:       resty.New(),
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
				url := a.getPathToMetric(metricType.String(), metricName.String(), metricValue)
				_, err := a.client.R().
					SetHeader("Content-Type", "text/plain").
					Post(url)
				if err != nil {
					log.Println(err)
					continue
				}

				//resp.Body.Close()
				//log.Println("response status:", resp.Status)
			}
		}
	}
}

func (a *Agent) Stop() {

}

func (a *Agent) getPathToMetric(metricType, metricName, metricValue string) string {
	return strings.Join([]string{"http:/", a.addr, "value", metricType, metricName, metricValue}, "/")
}

func (a *Agent) GetPollCounter() int64 {
	return a.pollCounter
}
