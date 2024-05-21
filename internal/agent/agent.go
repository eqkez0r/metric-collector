package agent

import (
	"context"
	"errors"
	"github.com/Eqke/metric-collector/internal/agent/generator"
	"github.com/Eqke/metric-collector/internal/agent/poller"
	"github.com/Eqke/metric-collector/internal/agent/poster"
	"github.com/Eqke/metric-collector/internal/agent/reqtype"
	"github.com/Eqke/metric-collector/internal/config"
	"github.com/Eqke/metric-collector/pkg/metric"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

const (
	errPointPostMetrics = "error in agent.postMetrics(): "
)

var (
	ErrUnknownMetricType = errors.New("unknown metric type")
)

type MetricPoller interface {
	Poll() metric.MetricMap
}

type MetricGenerator interface {
	Generate(mp metric.MetricMap) chan *reqtype.ReqType
	Shutdown()
}

type MetricPoster interface {
	Post(requests <-chan *reqtype.ReqType)
}

type Agent struct {
	logger      *zap.SugaredLogger
	settings    *config.AgentConfig
	client      *resty.Client
	pollCounter int64
	mp          metric.MetricMap
	mu          sync.RWMutex

	poller    MetricPoller
	generator MetricGenerator
	poster    MetricPoster
}

func New(
	settings *config.AgentConfig,
	logger *zap.SugaredLogger) *Agent {
	client := resty.New()

	return &Agent{
		logger:      logger,
		settings:    settings,
		client:      client,
		pollCounter: 0,
		mp:          make(metric.MetricMap),
		poller:      poller.NewPoller(logger),
		generator:   generator.NewGenerator(logger, settings),
		poster:      poster.NewPoster(logger, settings),
		mu:          sync.RWMutex{},
	}
}

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
		}
	}
}

func (a *Agent) updCounter() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.pollCounter++
	a.mp[metric.TypeCounter][metric.PollCount] = strconv.FormatInt(a.pollCounter, 10)
}
