// Пакет agent содержит агента(клиента), который отправляет запросы на
// сервер хранения метрик
package httpagent

import (
	"context"
	"crypto/rsa"
	"github.com/Eqke/metric-collector/internal/agent/config"
	"strconv"
	"sync"
	"time"

	"github.com/Eqke/metric-collector/internal/agent/generator"
	"github.com/Eqke/metric-collector/internal/agent/poller"
	"github.com/Eqke/metric-collector/internal/agent/poster"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// Структура Agent, который отправляет запросы на сервер
type Agent struct {
	logger      *zap.SugaredLogger
	client      *resty.Client
	pollCounter int64
	mp          metric.Map
	mu          sync.RWMutex
	poller      poller.MetricPoller
	generator   generator.MetricGenerator
	poster      poster.MetricPoster

	pollInterval   time.Duration
	reportInterval time.Duration
}

// Функция New возвращает объект агента
func New(
	settings *config.AgentConfig,
	logger *zap.SugaredLogger,
	publicKey *rsa.PublicKey,
	poller poller.MetricPoller,
) *Agent {
	client := resty.New()

	return &Agent{
		logger:         logger,
		client:         client,
		pollCounter:    0,
		mp:             make(metric.Map),
		poller:         poller,
		generator:      generator.NewGenerator(logger, settings, publicKey),
		poster:         poster.NewPoster(logger, settings),
		mu:             sync.RWMutex{},
		pollInterval:   time.Duration(settings.PollInterval) * time.Second,
		reportInterval: time.Duration(settings.ReportInterval) * time.Second,
	}
}

// Функция Run запускает процесс сбора метрик и их отправку на сервер
func (a *Agent) Run(ctx context.Context, wg *sync.WaitGroup) {
	reportTicker := time.NewTicker(a.reportInterval)
	defer reportTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			{
				a.logger.Info("agent was stopped")
				wg.Done()
				return
			}
		case <-reportTicker.C:
			{
				a.logger.Info("posting...")

				a.mu.Lock()
				a.mp = a.poller.GetMetrics()
				a.updCounter()
				a.poster.Post(a.generator.Generate(a.mp))
				a.mu.Unlock()
				a.logger.Info("posting... done")
			}
		default:
			{

			}
		}
	}
}

// Функция updCounter инкрементирует счетчик
func (a *Agent) updCounter() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.pollCounter++
	a.mp[metric.TypeCounter][metric.PollCount] = strconv.FormatInt(a.pollCounter, 10)
}
