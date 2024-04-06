package agent

import (
	"context"
	"github.com/Eqke/metric-collector/internal/config"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/go-resty/resty/v2"
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	pollCounterName = "pollCounter"
	randName        = "random"
	typeGauge       = "gauge"
	typeCounter     = "counter"
)

type Agent struct {
	addr           string
	settings       *config.AgentConfig
	client         *resty.Client
	reportInterval time.Duration
	pollInterval   time.Duration
	pollCounter    int64
	mp             map[metric.TypeMetric]map[metric.MetricName]string
	ctx            context.Context
	wg             sync.WaitGroup
	mu             *sync.Mutex
}

func New(ctx context.Context, settings *config.AgentConfig) *Agent {
	return &Agent{
		ctx:            ctx,
		addr:           settings.AgentEndpoint,
		settings:       settings,
		client:         resty.New(),
		pollInterval:   time.Duration(settings.PollInterval) * time.Second,
		reportInterval: time.Duration(settings.ReportInterval) * time.Second,
		pollCounter:    0,
	}
}

func (a *Agent) Run() {
	log.Println("agent was started")
	a.mu = &sync.Mutex{}
	ms := &runtime.MemStats{}
	a.mp = metric.PrepareMetrics(ms)
	a.wg.Add(1)
	go a.pollMetric(ms)
	a.wg.Add(1)
	go a.postMetrics()
	a.wg.Wait()
	log.Println("agent was stopped")
}

func (a *Agent) postMetrics() {
	defer a.wg.Done()
	for range time.Tick(a.reportInterval) {
		select {
		case <-a.ctx.Done():
			{
				log.Println("post metrics was stopped")
				return
			}
		default:
			{
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
		}
	}
}

func (a *Agent) pollMetric(ms *runtime.MemStats) {
	defer a.wg.Done()
	for range time.Tick(a.pollInterval) {
		select {
		case <-a.ctx.Done():
			{
				log.Println("update metrics was stopped")
				return
			}
		default:
			{
				// Обертка сделана для того, чтобы можно было корректно сбросить mutex
				func() {
					a.mu.Lock()
					defer log.Println("update metrics")
					defer a.mu.Unlock()
					metric.UpdateMetrics(ms, a.mp)

				}()
			}
		}
	}
}

func (a *Agent) getPathToMetric(metricType, metricName, metricValue string) string {
	return strings.Join([]string{"http:/", a.addr, "update", metricType, metricName, metricValue}, "/")
}
