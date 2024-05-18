package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/Eqke/metric-collector/internal/config"
	e "github.com/Eqke/metric-collector/pkg/error"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/Eqke/metric-collector/utils/hash"
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

	updatesEndpoint = "/updates"
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
	wg          sync.WaitGroup
	mu          *sync.Mutex
	attempt     int
}

func New(
	settings *config.AgentConfig,
	logger *zap.SugaredLogger) *Agent {
	client := resty.New()
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryAfter(func(c *resty.Client, r *resty.Response) (time.Duration, error) {
		return c.RetryWaitTime + 2*time.Duration(r.Request.Attempt)*time.Second, nil
	})
	client.SetRetryMaxWaitTime(5 * time.Second)
	return &Agent{
		logger:      logger,
		settings:    settings,
		client:      client,
		pollCounter: 0,
	}
}

func (a *Agent) Run(ctx context.Context) {
	a.mu = &sync.Mutex{}
	ms := &runtime.MemStats{}
	a.mp = metric.PrepareMetrics(ms)
	a.wg.Add(1)
	go a.pollMetricProcess(ctx, ms)
	a.wg.Add(1)
	go a.postMetrics(ctx)
	a.logger.Info("agent was started.")
	a.wg.Wait()
	a.logger.Info("agent was stopped")
}

func (a *Agent) postMetrics(ctx context.Context) {
	defer a.wg.Done()
	t := time.NewTicker(time.Duration(a.settings.ReportInterval) * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			{
				a.logger.Info("post metrics was stopped")
				return
			}
		case <-t.C:
			{
				a.updCounter()

				if err := a.pollMetricByBatch(); err != nil {
					a.logger.Errorf("%s: %v", errPointPostMetrics, err)
				}
				if err := a.pollEncodedMetricByBatch(); err != nil {
					a.logger.Errorf("%s: %v", errPointPostMetrics, err)
				}
			}
		}
	}
}

func (a *Agent) updCounter() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.pollCounter++
	a.mp[typeCounter][pollCounterName] = strconv.FormatInt(a.pollCounter, 10)
}

func (a *Agent) pollSingleMetric() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	for metricType, metricMap := range a.mp {
		for metricName, metricValue := range metricMap {
			a.logger.Infof("sending metric with type: %s, name: %s, value: %s",
				metricType, metricName, metricValue)
			if err := a.pollUsualMetric(metricName.String(), metricType.String(), metricValue); err != nil {
				a.logger.Errorf("%s: %v", errPointPostMetrics, err)
			}
			if err := a.pollJSONMetric(metricName.String(), metricType.String(), metricValue); err != nil {
				a.logger.Errorf("%s: %v", errPointPostMetrics, err)
			}
			if err := a.pollEncodeMetric(metricName.String(), metricType.String(), metricValue); err != nil {
				a.logger.Errorf("%s: %v", errPointPostMetrics, err)
			}
		}
	}
	return nil
}

func (a *Agent) pollUsualMetric(metricName, metricType, metricValue string) error {
	endPoint := a.getEndpointToUsualMetric(metricType, metricName, metricValue)
	var resp *resty.Response
	var err error
	if a.settings.HashKey == "" {
		resp, err = a.client.R().
			SetHeader("Content-Type", "text/plain").
			Post(endPoint)
	} else {
		h := hash.Hash([]byte{}, a.settings.HashKey)
		sign := base64.StdEncoding.EncodeToString(h)
		a.logger.Infof("hash: %s", sign)
		resp, err = a.client.R().
			SetHeader("Content-Type", "text/plain").
			SetHeader("HashSHA256", sign).
			Post(endPoint)
	}
	if err != nil {
		log.Println(e.WrapError(errPointPostMetrics, err))
		return err
	}
	a.logger.Infof("endpoint: %s, status_code: %d, size: %d",
		endPoint, resp.StatusCode(), resp.Size())
	return nil
}

func (a *Agent) pollJSONMetric(metricName, metricType, metricValue string) error {
	b, err := a.prepareJSONMetric(metricName, metricType, metricValue)
	if err != nil {
		a.logger.Errorf("%s: %v", errPointPostMetrics, err)
		return err
	}
	endPoint := a.getEndpointToJSONMetric()
	var resp *resty.Response
	if a.settings.HashKey == "" {
		resp, err = a.client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(b).
			Post(endPoint)
	} else {
		h := hash.Hash(b, a.settings.HashKey)
		sign := base64.StdEncoding.EncodeToString(h)
		a.logger.Infof("hash: %s", sign)
		resp, err = a.client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("HashSHA256", sign).
			SetBody(b).
			Post(endPoint)
	}
	if err != nil {
		a.logger.Errorf("%s: %v", errPointPostMetrics, err)
		return err
	}
	a.logger.Infof("endpoint: %s, status_code: %d, size: %d",
		endPoint, resp.StatusCode(), resp.Size())
	return nil
}

func (a *Agent) pollEncodeMetric(metricName, metricType, metricValue string) error {
	b, err := a.prepareJSONMetric(metricName, metricType, metricValue)
	if err != nil {
		a.logger.Errorf("%s: %v", errPointPostMetrics, err)
		return err
	}

	endPoint := a.getEndpointToJSONMetric()
	encoded, err := a.compress(b)
	if err != nil {
		a.logger.Errorf("%s: %v", errPointPostMetrics, err)
		return err
	}

	var resp *resty.Response
	if a.settings.HashKey == "" {
		resp, err = a.client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(encoded).
			Post(endPoint)
	} else {
		h := hash.Hash(encoded, a.settings.HashKey)
		sign := base64.StdEncoding.EncodeToString(h)
		a.logger.Infof("hash: %s", sign)
		resp, err = a.client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("HashSHA256", sign).
			SetBody(encoded).
			Post(endPoint)
	}

	if err != nil {
		a.logger.Errorf("%s: %v", errPointPostMetrics, err)
		return err
	}
	a.logger.Infof("endpoint: %s, status_code: %d, size: %d",
		endPoint, resp.StatusCode(), resp.Size())
	return nil
}

func (a *Agent) pollMetricByBatch() error {
	a.mu.Lock()
	arr := []metric.Metrics{}
	for k, v := range a.mp {
		switch k {
		case metric.TypeGauge:
			{
				for metricName, metricValue := range v {
					val, err := strconv.ParseFloat(metricValue, 64)
					if err != nil {
						a.logger.Errorf("%s: %v", errPointPostMetrics, err)
						continue
					}
					met := metric.Metrics{ID: metricName.String(), MType: k.String(), Value: &val}
					arr = append(arr, met)
				}
			}
		case metric.TypeCounter:
			{
				for metricName, metricValue := range v {
					val, err := strconv.ParseInt(metricValue, 10, 64)
					if err != nil {
						a.logger.Errorf("%s: %v", errPointPostMetrics, err)
						continue
					}
					met := metric.Metrics{ID: metricName.String(), MType: k.String(), Delta: &val}
					arr = append(arr, met)
				}
			}
		}
	}
	a.mu.Unlock()
	if len(arr) == 0 {
		return nil
	}
	var resp *resty.Response
	var err error
	var b []byte
	if a.settings.HashKey == "" {
		resp, err = a.client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(arr).
			Post(a.getEndpointToBatchMetric())

	} else {
		b, err = json.Marshal(arr)
		if err != nil {
			a.logger.Errorf("%s: %v", errPointPostMetrics, err)
			return err
		}
		h := hash.Hash(b, a.settings.HashKey)
		sign := base64.StdEncoding.EncodeToString(h)
		a.logger.Infof("hash: %s", sign)
		//a.logger.Infof("data: %s", b)
		resp, err = a.client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("HashSHA256", sign).
			SetBody(arr).
			Post(a.getEndpointToBatchMetric())

	}
	if err != nil {
		a.logger.Errorf("%s: %v", errPointPostMetrics, err)
		return err
	}
	a.logger.Infof("endpoint: %s, status_code: %d, size: %d",
		a.getEndpointToBatchMetric(), resp.StatusCode(), resp.Size())
	return nil
}

func (a *Agent) pollEncodedMetricByBatch() error {
	a.mu.Lock()
	arr := []metric.Metrics{}
	for k, v := range a.mp {
		switch k {
		case metric.TypeGauge:
			{
				for metricName, metricValue := range v {
					val, err := strconv.ParseFloat(metricValue, 64)
					if err != nil {
						a.logger.Errorf("%s: %v", errPointPostMetrics, err)
						continue
					}
					met := metric.Metrics{ID: metricName.String(), MType: k.String(), Value: &val}
					arr = append(arr, met)
				}
			}
		case metric.TypeCounter:
			{
				for metricName, metricValue := range v {
					val, err := strconv.ParseInt(metricValue, 10, 64)
					if err != nil {
						a.logger.Errorf("%s: %v", errPointPostMetrics, err)
						continue
					}
					met := metric.Metrics{ID: metricName.String(), MType: k.String(), Delta: &val}
					arr = append(arr, met)
				}
			}
		}
	}
	a.mu.Unlock()
	if len(arr) == 0 {
		return nil
	}
	b, err := json.Marshal(arr)
	if err != nil {
		a.logger.Errorf("%s: %v", errPointPostMetrics, err)
		return err
	}
	encoded, err := a.compress(b)
	if err != nil {
		a.logger.Errorf("%s: %v", errPointPostMetrics, err)
		return err
	}
	var resp *resty.Response
	if a.settings.HashKey == "" {
		resp, err = a.client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetBody(encoded).
			Post(a.getEndpointToBatchMetric())
	} else {
		h := hash.Hash(encoded, a.settings.HashKey)
		sign := base64.StdEncoding.EncodeToString(h)
		a.logger.Infof("hash: %s", sign)
		//a.logger.Infof("data: %s", encoded)
		resp, err = a.client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("HashSHA256", sign).
			SetHeader("Content-Encoding", "gzip").
			SetBody(encoded).
			Post(a.getEndpointToBatchMetric())
	}
	if err != nil {
		a.logger.Errorf("%s: %v", errPointPostMetrics, err)
		return err
	}
	a.logger.Infof("endpoint: %s, status_code: %d, size: %d",
		a.getEndpointToBatchMetric(), resp.StatusCode(), resp.Size())
	return nil
}

func (a *Agent) pollMetricProcess(ctx context.Context, ms *runtime.MemStats) {
	defer a.wg.Done()
	t := time.NewTicker(time.Duration(a.settings.PollInterval) * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			{
				a.logger.Info("poll metrics was stopped")
				return
			}
		case <-t.C:
			{
				a.pollMetric(ms)
			}
		}
	}
}

func (a *Agent) pollMetric(ms *runtime.MemStats) {
	a.mu.Lock()
	defer a.logger.Info("update metrics")
	defer a.mu.Unlock()
	metric.UpdateMetrics(ms, a.mp)
}

func (a *Agent) getEndpointToUsualMetric(metricType, metricName, metricValue string) string {
	return strings.Join([]string{"http:/", a.settings.AgentEndpoint, "update", metricType, metricName, metricValue}, "/")
}

func (a *Agent) getEndpointToJSONMetric() string {
	return strings.Join([]string{"http:/", a.settings.AgentEndpoint, "update"}, "/")
}

func (a *Agent) getEndpointToBatchMetric() string {
	return strings.Join([]string{"http:/", a.settings.AgentEndpoint, "updates"}, "/")
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

	b, err := json.Marshal(m)
	if err != nil {
		a.logger.Errorf("%s: %v", errPointPostMetrics, err)
		return nil, err
	}
	return b, nil
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
