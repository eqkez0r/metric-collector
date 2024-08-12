package generator

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"github.com/Eqke/metric-collector/internal/agent/config"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Eqke/metric-collector/internal/agent/reqtype"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/Eqke/metric-collector/utils/hash"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

const (
	errPointPostMetrics = "error in generator.postMetrics(): "
)

var (
	ErrEmptyMetricBatch  = errors.New("empty batch")
	ErrUnknownMetricType = errors.New("unknown metric type")
)

type MetricGenerator interface {
	Generate(mp metric.Map) chan *reqtype.ReqType
	Shutdown()
}

type Generator struct {
	logger            *zap.SugaredLogger
	mu                sync.Mutex
	settings          *config.AgentConfig
	generatedRequests chan *reqtype.ReqType
	mp                metric.Map
	errChan           chan error
	client            *resty.Client
}

func NewGenerator(
	logger *zap.SugaredLogger,
	settings *config.AgentConfig) *Generator {
	client := resty.New()
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryAfter(func(c *resty.Client, r *resty.Response) (time.Duration, error) {
		return c.RetryWaitTime + 2*time.Duration(r.Request.Attempt)*time.Second, nil
	})
	client.SetRetryMaxWaitTime(5 * time.Second)
	return &Generator{
		logger:            logger,
		settings:          settings,
		generatedRequests: make(chan *reqtype.ReqType, settings.RateLimit),
		mu:                sync.Mutex{},
		errChan:           make(chan error),
		client:            client,
	}
}

func (g *Generator) Generate(mp metric.Map) chan *reqtype.ReqType {
	g.mp = mp
	done := make(chan struct{})

	go g.errorLogger(done)

	go g.pollSingleMetric()
	go g.pollMetricByBatch()
	go g.pollEncodedMetricByBatch()

	close(done)
	return g.generatedRequests
}

func (g *Generator) Shutdown() {
	close(g.generatedRequests)
	close(g.errChan)
}

func (g *Generator) errorLogger(done chan struct{}) {
	for {
		select {
		case <-done:
			g.logger.Info("generator was stopped")
			return
		case err := <-g.errChan:
			g.logger.Error(err)
		default:
			{

			}
		}
	}
}

func (g *Generator) pollSingleMetric() {
	g.mu.Lock()
	m := make(metric.Map)
	m[metric.TypeCounter] = make(map[metric.Name]string)
	m[metric.TypeGauge] = make(map[metric.Name]string)
	for k, v := range g.mp[metric.TypeCounter] {
		m[metric.TypeCounter][k] = v
	}
	for k, v := range g.mp[metric.TypeGauge] {
		m[metric.TypeGauge][k] = v
	}
	g.mu.Unlock()
	for metricType, metricMap := range m {
		for metricName, metricValue := range metricMap {
			g.logger.Infof("sending metric with type: %s, name: %s, value: %s",
				metricType, metricName, metricValue)
			req, err := g.pollUsualMetric(metricName.String(), metricType.String(), metricValue)
			if err != nil {
				g.errChan <- err
			} else {
				g.generatedRequests <- req
			}
			req, err = g.pollJSONMetric(metricName.String(), metricType.String(), metricValue)
			if err != nil {
				g.errChan <- err
			} else {
				g.generatedRequests <- req
			}
			req, err = g.pollEncodeMetric(metricName.String(), metricType.String(), metricValue)
			if err != nil {
				g.errChan <- err
			} else {
				g.generatedRequests <- req
			}
		}
	}
}

func (g *Generator) pollUsualMetric(metricName, metricType, metricValue string) (*reqtype.ReqType, error) {
	endPoint := g.getEndpointToUsualMetric(metricType, metricName, metricValue)
	req := g.client.R().SetHeader("Content-Type", "text/plain")
	return &reqtype.ReqType{Req: req, Endpoint: endPoint}, nil
}

func (g *Generator) pollJSONMetric(metricName, metricType, metricValue string) (*reqtype.ReqType, error) {
	b, err := g.prepareJSONMetric(metricName, metricType, metricValue)
	if err != nil {
		return nil, err
	}
	endPoint := g.getEndpointToJSONMetric()
	req := g.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(b)
	if g.settings.HashKey != "" {
		req = req.SetHeader("HashSHA256", hash.Sign(b, g.settings.HashKey))
	}

	return &reqtype.ReqType{Req: req, Endpoint: endPoint}, nil
}

func (g *Generator) pollEncodeMetric(metricName, metricType, metricValue string) (*reqtype.ReqType, error) {
	b, err := g.prepareJSONMetric(metricName, metricType, metricValue)
	if err != nil {
		return nil, err
	}

	endPoint := g.getEndpointToJSONMetric()
	encoded, err := g.compress(b)
	if err != nil {
		return nil, err
	}

	req := g.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(encoded)
	if g.settings.HashKey != "" {
		req = req.SetHeader("HashSHA256", hash.Sign(encoded, g.settings.HashKey))
	}
	return &reqtype.ReqType{Req: req, Endpoint: endPoint}, nil
}

func (g *Generator) pollMetricByBatch() {
	g.mu.Lock()
	arr := g.prepareMetricArray()
	g.mu.Unlock()
	if len(arr) == 0 {
		g.errChan <- ErrEmptyMetricBatch
		return
	}
	req := g.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(arr)
	if g.settings.HashKey != "" {
		b, err := json.Marshal(arr)
		if err != nil {
			g.errChan <- err
			return
		}
		req = req.SetHeader("HashSHA256", hash.Sign(b, g.settings.HashKey))
	}
	endpoint := g.getEndpointToBatchMetric()
	g.generatedRequests <- &reqtype.ReqType{Req: req, Endpoint: endpoint}
}

func (g *Generator) pollEncodedMetricByBatch() {
	g.mu.Lock()
	arr := g.prepareMetricArray()
	g.mu.Unlock()
	if len(arr) == 0 {
		g.errChan <- ErrEmptyMetricBatch
		return
	}
	b, err := json.Marshal(arr)
	if err != nil {
		g.errChan <- err
		return
	}
	encoded, err := g.compress(b)
	if err != nil {
		g.errChan <- err
		return
	}
	req := g.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(encoded)
	if g.settings.HashKey != "" {
		req = req.SetHeader("HashSHA256", hash.Sign(encoded, g.settings.HashKey))
	}
	g.generatedRequests <- &reqtype.ReqType{Req: req, Endpoint: g.getEndpointToBatchMetric()}
}

func (g *Generator) getEndpointToUsualMetric(metricType, metricName, metricValue string) string {
	return strings.Join([]string{"http:/", g.settings.AgentEndpoint, "update", metricType, metricName, metricValue}, "/")
}

func (g *Generator) getEndpointToJSONMetric() string {
	return strings.Join([]string{"http:/", g.settings.AgentEndpoint, "update"}, "/")
}

func (g *Generator) getEndpointToBatchMetric() string {
	return strings.Join([]string{"http:/", g.settings.AgentEndpoint, "updates"}, "/")
}

func (g *Generator) prepareJSONMetric(metricName, metricType, metricValue string) ([]byte, error) {
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
				g.logger.Errorf("%s: %v", errPointPostMetrics, err)
				return nil, err
			}
			m.Value = &val
		}
	case metric.TypeCounter.String():
		{
			val, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				g.logger.Errorf("%s: %v", errPointPostMetrics, err)
				return nil, err
			}
			m.Delta = &val
		}
	default:
		{
			g.logger.Errorf("unknown metric type: %s", metricType)
			return nil, ErrUnknownMetricType
		}
	}

	b, err := json.Marshal(m)
	if err != nil {
		g.logger.Errorf("%s: %v", errPointPostMetrics, err)
		return nil, err
	}
	return b, nil
}

func (g *Generator) prepareMetricArray() []metric.Metrics {
	arr := make([]metric.Metrics, 0, len(g.mp))
	for k, v := range g.mp {
		switch k {
		case metric.TypeGauge:
			{
				for metricName, metricValue := range v {
					val, err := strconv.ParseFloat(metricValue, 64)
					if err != nil {
						g.errChan <- err
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
						g.errChan <- err
						continue
					}
					met := metric.Metrics{ID: metricName.String(), MType: k.String(), Delta: &val}
					arr = append(arr, met)
				}
			}
		}
	}
	return arr
}

func (g *Generator) compress(b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(buf)

	_, err := gz.Write(b)
	if err != nil {
		return nil, err
	}
	gz.Close()

	return buf.Bytes(), nil
}

func (g *Generator) decompress(b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(b)
	gz, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(gz)
}
