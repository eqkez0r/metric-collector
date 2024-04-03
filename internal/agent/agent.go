package agent

import (
	"context"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/go-resty/resty/v2"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	pollCounterName = "pollCounter"
	randName        = "random"
	typeGauge       = "gauge"
	typeCounter     = "counter"
)

type Agent struct {
	mu             sync.RWMutex
	addr           string
	client         *resty.Client
	reportInterval time.Duration
	pollInterval   time.Duration
	pollCounter    int64
	mp             map[metric.TypeMetric]map[metric.MetricName]string
	ctx            context.Context
}

func New(ctx context.Context, addr string, pollInterval, reportInterval int) *Agent {
	return &Agent{
		addr:           addr,
		client:         resty.New(),
		pollInterval:   time.Duration(pollInterval) * time.Second,
		reportInterval: time.Duration(reportInterval) * time.Second,
		pollCounter:    0,
		ctx:            ctx,
	}
}

func (a *Agent) Run() {
	log.Println("agent was started")
	signal.NotifyContext(a.ctx, syscall.SIGINT, syscall.SIGTERM)
	ms := &runtime.MemStats{}
	a.mp = metric.PrepareMetrics(ms)
	go a.updateMetrics(ms)
	go func() {
		for range time.Tick(a.reportInterval) {
			a.pollCounter++
			a.mp[typeCounter][pollCounterName] = strconv.FormatInt(a.pollCounter, 10)
			for metricType, metricMap := range a.mp {
				for metricName, metricValue := range metricMap {
					endPoint := a.getPathToMetric(metricType.String(), metricName.String(), metricValue)
					resp, err := a.client.R().
						SetHeader("Content-Type", "text/plain").
						Post(endPoint)
					if err != nil {
						log.Println(err)
						continue
					}

					log.Println(metricName, "was updated. response status code:", resp.StatusCode())
				}
			}
		}
	}()
	<-a.ctx.Done()
	a.Shutdown()
}

func (a *Agent) updateMetrics(ms *runtime.MemStats) {
	for range time.Tick(a.pollInterval) {
		a.mu.Lock()
		metric.UpdateMetrics(ms, a.mp)
		a.mu.Unlock()
		log.Println("metrics were updated")
	}
}

func (a *Agent) Shutdown() {
	log.Println("agent was stopped")
	os.Exit(1)
}

func (a *Agent) getPathToMetric(metricType, metricName, metricValue string) string {
	return strings.Join([]string{"http:/", a.addr, "update", metricType, metricName, metricValue}, "/")
}
