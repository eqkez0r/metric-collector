package agent

import (
	"context"
	"encoding/json"
	"github.com/Eqke/metric-collector/internal/config"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	pollCounterName = "PollCount"
	randName        = "random"
	typeGauge       = "gauge"
	typeCounter     = "counter"

	errPointPostMetrics = "error in agent.postMetrics(): "
)

type Agent struct {
	logger      *zap.SugaredLogger
	settings    *config.AgentConfig
	client      *resty.Client
	pollCounter int64
	mp          map[metric.TypeMetric]map[metric.MetricName]string
	ctx         context.Context
	wg          sync.WaitGroup
	mu          *sync.Mutex
}

func New(
	ctx context.Context,
	settings *config.AgentConfig,
	logger *zap.SugaredLogger) *Agent {
	return &Agent{
		ctx:         ctx,
		logger:      logger,
		settings:    settings,
		client:      resty.New(),
		pollCounter: 0,
	}
}

func (a *Agent) Run() {
	a.mu = &sync.Mutex{}
	ms := &runtime.MemStats{}
	a.mp = metric.PrepareMetrics(ms)
	a.wg.Add(1)
	go a.pollMetric(ms)
	a.wg.Add(1)
	go a.postMetrics()
	a.logger.Info("agent was started.")
	a.wg.Wait()
	a.logger.Info("agent was stopped")
}

func (a *Agent) postMetrics() {
	defer a.wg.Done()
	ticker := time.NewTicker(time.Duration(a.settings.ReportInterval) * time.Second)
	for {
		select {
		case <-a.ctx.Done():
			{
				a.logger.Info("post metrics was stopped")
				return
			}
		case <-ticker.C:
			{
				a.pollCounter++
				a.mp[typeCounter][pollCounterName] = strconv.FormatInt(a.pollCounter, 10)

				for metricType, metricMap := range a.mp {
					for metricName, metricValue := range metricMap {
						a.logger.Infof("sending metric with type: %s, name: %s, value: %s",
							metricType, metricName, metricValue)
						m := metric.Metrics{
							ID:    metricName.String(),
							MType: metricType.String(),
							Delta: nil,
							Value: nil,
						}
						switch metricType {
						case metric.TypeGauge:
							{
								val, err := strconv.ParseFloat(metricValue, 64)
								if err != nil {
									a.logger.Errorf("%s: %v", errPointPostMetrics, err)
									continue
								}
								m.Value = &val
							}
						case metric.TypeCounter:
							{
								val, err := strconv.ParseInt(metricValue, 10, 64)
								if err != nil {
									a.logger.Errorf("%s: %v", errPointPostMetrics, err)
									continue
								}
								m.Delta = &val
							}
						default:
							{
								a.logger.Errorf("unknown metric type: %s", metricType)
								continue
							}
						}
						bytes, err := json.Marshal(m)
						if err != nil {
							a.logger.Errorf("%s: %v", errPointPostMetrics, err)
							continue
						}
						endPoint := "http://" + a.settings.AgentEndpoint + "/update"
						resp, err := a.client.R().
							SetHeader("Content-Type", "application/json").
							SetBody(bytes).Post(endPoint)
						if err != nil {
							a.logger.Errorf("%s: %v", errPointPostMetrics, err)
							continue
						}
						a.logger.Infof("endpoint: %s, status_code: %d, size: %d",
							endPoint, resp.StatusCode(), resp.Size())
						//endPoint := a.getPathToMetric(metricType.String(), metricName.String(), metricValue)
						//resp, err := a.client.R().
						//	SetHeader("Content-Type", "text/plain").
						//	Post(endPoint)
						//if err != nil {
						//	a.logger.Errorf("%s: %v", errPointPostMetrics, err)
						//	continue
						//}
						//a.logger.Infof("endpoint: %s, status_code: %d, size: %d",
						//	endPoint, resp.StatusCode(), resp.Size())
					}
				}
			}
		}
	}
}

func (a *Agent) pollMetric(ms *runtime.MemStats) {
	defer a.wg.Done()
	ticker := time.NewTicker(time.Duration(a.settings.PollInterval) * time.Second)
	for {
		select {
		case <-a.ctx.Done():
			{
				a.logger.Info("poll metrics was stopped")
				return
			}
		case <-ticker.C:
			{
				// Обертка сделана для того, чтобы можно было корректно сбросить mutex
				func() {
					a.mu.Lock()
					defer a.logger.Info("update metrics")
					defer a.mu.Unlock()
					metric.UpdateMetrics(ms, a.mp)
				}()

			}
		}
	}
}

func (a *Agent) getPathToMetric(metricType, metricName, metricValue string) string {
	return strings.Join([]string{"http:/", a.settings.AgentEndpoint, "update", metricType, metricName, metricValue}, "/")
}
