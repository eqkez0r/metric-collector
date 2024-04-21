package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"github.com/Eqke/metric-collector/internal/config"
	e "github.com/Eqke/metric-collector/pkg/error"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"log"
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

var (
	ErrUnknownMetricType = errors.New("unknown metric type")
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
						//if err := a.pollUsualMetric(metricName.String(), metricType.String(), metricValue); err != nil {
						//	a.logger.Errorf("%s: %v", errPointPostMetrics, err)
						//}
						if err := a.pollJSONMetric(metricName.String(), metricType.String(), metricValue); err != nil {
							a.logger.Errorf("%s: %v", errPointPostMetrics, err)
						}
						time.Sleep(time.Millisecond * 500)
						if err := a.pollEncodeMetric(metricName.String(), metricType.String(), metricValue); err != nil {
							a.logger.Errorf("%s: %v", errPointPostMetrics, err)
						}
					}
				}
			}
		}
	}
}

func (a *Agent) pollUsualMetric(metricName, metricType, metricValue string) error {
	endPoint := a.getEndpointToUsualMetric(metricType, metricName, metricValue)
	resp, err := a.client.R().
		SetHeader("Content-Type", "text/plain").
		Post(endPoint)
	if err != nil {
		log.Println(e.WrapError(errPointPostMetrics, err))
		return err
	}
	a.logger.Infof("endpoint: %s, status_code: %d, size: %d",
		endPoint, resp.StatusCode(), resp.Size())
	return nil
}

func (a *Agent) pollJSONMetric(metricName, metricType, metricValue string) error {
	bytes, err := a.prepareJSONMetric(metricName, metricType, metricValue)
	if err != nil {
		a.logger.Errorf("%s: %v", errPointPostMetrics, err)
		return err
	}
	endPoint := a.getEndpointToJSONMetric()
	resp, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bytes).
		Post(endPoint)
	if err != nil {
		a.logger.Errorf("%s: %v", errPointPostMetrics, err)
		return err
	}
	a.logger.Infof("endpoint: %s, status_code: %d, size: %d",
		endPoint, resp.StatusCode(), resp.Size())
	return nil
}

func (a *Agent) pollEncodeMetric(metricName, metricType, metricValue string) error {
	bytes, err := a.prepareJSONMetric(metricName, metricType, metricValue)
	if err != nil {
		a.logger.Errorf("%s: %v", errPointPostMetrics, err)
		return err
	}

	endPoint := a.getEndpointToJSONMetric()
	encoded, err := a.compress(bytes)
	if err != nil {
		a.logger.Errorf("%s: %v", errPointPostMetrics, err)
		return err
	}

	resp, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(encoded).
		Post(endPoint)

	if err != nil {
		a.logger.Errorf("%s: %v", errPointPostMetrics, err)
		return err
	}
	a.logger.Infof("endpoint: %s, status_code: %d, size: %d",
		endPoint, resp.StatusCode(), resp.Size())
	return nil
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

func (a *Agent) getEndpointToUsualMetric(metricType, metricName, metricValue string) string {
	return strings.Join([]string{"http:/", a.settings.AgentEndpoint, "update", metricType, metricName, metricValue}, "/")
}

func (a *Agent) getEndpointToJSONMetric() string {
	return strings.Join([]string{"http:/", a.settings.AgentEndpoint, "update"}, "/")
}

func (a *Agent) prepareJSONMetric(metricName, metricType, metricValue string) ([]byte, error) {
	m := metric.Metrics{
		ID:    metricName,
		MType: metricType,
		Delta: nil,
		Value: nil,
	}

	switch metricType {
	case metric.TypeGauge.String():
		{
			val, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				a.logger.Errorf("%s: %v", errPointPostMetrics, err)
				return nil, err
			}
			m.Value = &val
		}
	case metric.TypeCounter.String():
		{
			val, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				a.logger.Errorf("%s: %v", errPointPostMetrics, err)
				return nil, err
			}
			m.Delta = &val
		}
	default:
		{
			a.logger.Errorf("unknown metric type: %s", metricType)
			return nil, ErrUnknownMetricType
		}
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		a.logger.Errorf("%s: %v", errPointPostMetrics, err)
		return nil, err
	}
	return bytes, nil
}

func (a *Agent) compress(b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(buf)

	_, err := gz.Write(b)
	if err != nil {
		return nil, err
	}
	gz.Close()

	return buf.Bytes(), nil
}
