// Пакет generator предоставляет работу с генератором запросов
package generator

import (
	"bytes"
	"compress/gzip"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"github.com/Eqke/metric-collector/internal/agent/config"
	"github.com/Eqke/metric-collector/internal/encrypting"
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

// Объявление
var (
	ErrEmptyMetricBatch  = errors.New("empty batch")
	ErrUnknownMetricType = errors.New("unknown metric type")
)

// Интерфейс MetricGenerator предоставляет генератор запросов
type MetricGenerator interface {
	Generate(mp metric.Map) chan *reqtype.ReqType
	Shutdown()
}

// Тип Generator является реализацией интерфейса MetricGenerator
type Generator struct {
	logger            *zap.SugaredLogger
	mu                sync.Mutex
	settings          *config.AgentConfig
	generatedRequests chan *reqtype.ReqType
	mp                metric.Map
	errChan           chan error
	client            *resty.Client
	publicKey         *rsa.PublicKey
}

// Функция NewGenerator возвращает экземпляр Generator
func NewGenerator(
	logger *zap.SugaredLogger,
	settings *config.AgentConfig,
	publicKey *rsa.PublicKey,
) *Generator {
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
		publicKey:         publicKey,
	}
}

// Метод Generate производит сбор метрик и возвращает канал с запросами
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

// Метод Shutdown позволяет корректно высвободить ресурсы
func (g *Generator) Shutdown() {
	close(g.generatedRequests)
	close(g.errChan)
}

// Метод errorLogger представляет собой обработчик ошибок
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

// Метод pollSingleMetric отвечает за публикацию метрик
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

// Метод pollUsualMetric отвечает за получение реквеста единичной метрики
func (g *Generator) pollUsualMetric(metricName, metricType, metricValue string) (*reqtype.ReqType, error) {
	endPoint := g.getEndpointToUsualMetric(metricType, metricName, metricValue)
	req := g.client.R().SetHeader("Content-Type", "text/plain")
	return &reqtype.ReqType{Req: req, Endpoint: endPoint}, nil
}

// Метод pollJSONMetric формирует запрос для метрики в формате JSON
func (g *Generator) pollJSONMetric(metricName, metricType, metricValue string) (*reqtype.ReqType, error) {
	b, err := g.prepareJSONMetric(metricName, metricType, metricValue)
	if err != nil {
		return nil, err
	}
	endPoint := g.getEndpointToJSONMetric()
	encryptedData, err := encrypting.Encrypt(g.publicKey, b)
	if err != nil {
		return nil, err
	}
	req := g.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(encryptedData)
	if g.settings.HashKey != "" {
		req = req.SetHeader("HashSHA256", hash.Sign(encryptedData, g.settings.HashKey))
	}

	return &reqtype.ReqType{Req: req, Endpoint: endPoint}, nil
}

// Метод pollEncodeMetric формирует шифрованный запрос
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

	encryptedData, err := encrypting.Encrypt(g.publicKey, encoded)
	if err != nil {
		return nil, err
	}

	req := g.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(encoded)
	if g.settings.HashKey != "" {
		req = req.SetHeader("HashSHA256", hash.Sign(encryptedData, g.settings.HashKey))
	}
	return &reqtype.ReqType{Req: req, Endpoint: endPoint}, nil
}

// Метод pollMetricByBatch формирует запрос пачкой метрик
func (g *Generator) pollMetricByBatch() {
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
	encryptedData, err := encrypting.Encrypt(g.publicKey, b)
	if err != nil {
		g.errChan <- err
		return
	}
	req := g.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(encryptedData)
	if g.settings.HashKey != "" {
		req = req.SetHeader("HashSHA256", hash.Sign(encryptedData, g.settings.HashKey))
	}
	endpoint := g.getEndpointToBatchMetric()
	g.generatedRequests <- &reqtype.ReqType{Req: req, Endpoint: endpoint}
}

// Метод pollEncodedMetricByBatch формирует запрос шифрованной пачки
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
	encryptedData, err := encrypting.Encrypt(g.publicKey, encoded)
	if err != nil {
		g.errChan <- err
		return
	}
	req := g.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(encryptedData)
	if g.settings.HashKey != "" {
		req = req.SetHeader("HashSHA256", hash.Sign(encryptedData, g.settings.HashKey))
	}
	g.generatedRequests <- &reqtype.ReqType{Req: req, Endpoint: g.getEndpointToBatchMetric()}
}

// Метод getEndpointToUsualMetric формирует конечную точку для запроса
func (g *Generator) getEndpointToUsualMetric(metricType, metricName, metricValue string) string {
	return strings.Join([]string{"http:/", g.settings.AgentEndpoint, "update", metricType, metricName, metricValue}, "/")
}

// Метод getEndpointToJSONMetric формирует конечную точку для запроса в формате JSON
func (g *Generator) getEndpointToJSONMetric() string {
	return strings.Join([]string{"http:/", g.settings.AgentEndpoint, "update"}, "/")
}

// Метод getEndpointToBatchMetric формирует конечную точку для запроса в формате пачки
func (g *Generator) getEndpointToBatchMetric() string {
	return strings.Join([]string{"http:/", g.settings.AgentEndpoint, "updates"}, "/")
}

// Метод prepareJSONMetric отвечает за подготовку метрики в формате JSON
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

// Метод prepareMetricArray отвечает за подготовку массива метрик
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

// Метод compress отвечает за сжатие
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

// Метод decompress отвечает за разжатие
func (g *Generator) decompress(b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(b)
	gz, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(gz)
}
