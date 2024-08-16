// Пакет agent содержит агента(клиента), который отправляет запросы на
// сервер хранения метрик
package agent

import (
	"context"
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
	settings    *config.AgentConfig
	client      *resty.Client
	pollCounter int64
	mp          metric.Map
	mu          sync.RWMutex

	poller    poller.MetricPoller
	generator generator.MetricGenerator
	poster    poster.MetricPoster
}

// Функция New возвращает объект агента
func New(
	settings *config.AgentConfig,
	logger *zap.SugaredLogger) *Agent {
	client := resty.New()

	return &Agent{
		logger:      logger,
		settings:    settings,
		client:      client,
		pollCounter: 0,
		mp:          make(metric.Map),
		poller:      poller.NewPoller(logger),
		generator:   generator.NewGenerator(logger, settings),
		poster:      poster.NewPoster(logger, settings),
		mu:          sync.RWMutex{},
	}
}

// Функция Run запускает процесс сбора метрик и их отправку на сервер
func (a *Agent) Run(ctx context.Context) {
	pollTicker := time.NewTicker(time.Duration(a.settings.PollInterval) * time.Second)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(time.Duration(a.settings.ReportInterval) * time.Second)
	defer reportTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			{
				a.logger.Info("agent was stopped")
				return
			}
		case <-pollTicker.C:
			{
				a.logger.Info("polling...")
				a.mu.RLock()
				a.mp = a.poller.Poll()
				a.mu.RUnlock()
				a.logger.Info("polling... done")
			}
		case <-reportTicker.C:
			{
				a.logger.Info("posting...")
				a.updCounter()
				a.mu.Lock()
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
