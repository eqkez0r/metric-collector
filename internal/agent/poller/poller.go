// Пакет poller отвечает за сбор метрик
package poller

import (
	"context"
	"github.com/Eqke/metric-collector/internal/agent/config"
	"runtime"
	"sync"
	"time"

	"github.com/Eqke/metric-collector/pkg/metric"
	"go.uber.org/zap"
)

// Интерфейс MetricPoller отвечает за получение метрик
type MetricPoller interface {
	Poll(ctx context.Context, wg *sync.WaitGroup)
	GetMetrics() metric.Map
}

// Тип Poller является реализацией интерфейса MetricPoller
type Poller struct {
	logger       *zap.SugaredLogger
	mp           metric.Map
	mu           sync.Mutex
	pollInterval time.Duration
}

// Функция NewPoller возвращает объект типа Poller
func NewPoller(
	logger *zap.SugaredLogger,
	settings *config.AgentConfig,
) *Poller {
	mp := make(metric.Map)
	mp[metric.TypeGauge] = make(map[metric.Name]string)
	mp[metric.TypeCounter] = make(map[metric.Name]string)
	return &Poller{
		logger:       logger,
		mp:           mp,
		mu:           sync.Mutex{},
		pollInterval: time.Duration(settings.PollInterval) * time.Second,
	}
}

// Метод Poll собирает метрики и возвращает объект metric.Map
func (p *Poller) Poll(ctx context.Context, wg *sync.WaitGroup) {
	pollerTicker := time.NewTicker(p.pollInterval)
	defer pollerTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case <-pollerTicker.C:
			p.logger.Info("polling...")
			wg.Add(2)
			go p.updateRuntime(wg)
			go p.updateUtil(wg)
			p.logger.Info("polling complete")
		}
	}

}

func (p *Poller) GetMetrics() metric.Map {
	cp := make(metric.Map, len(p.mp))
	p.mu.Lock()
	defer p.mu.Unlock()
	for k, v := range p.mp {
		cp[k] = v
	}

	return cp
}

// Метод updateRuntime обновляет метрики из пакета runtime
func (p *Poller) updateRuntime(wg *sync.WaitGroup) {
	defer wg.Done()
	ms := &runtime.MemStats{}
	runtime.ReadMemStats(ms)
	metric.UpdateRuntimeMetrics(ms, p.mp)
}

// Метод updateUtil обновляет метрики cpu, totalmemory и freememory
func (p *Poller) updateUtil(wg *sync.WaitGroup) {
	defer wg.Done()
	if err := metric.UpdateUtilMetrics(p.mp); err != nil {
		p.logger.Errorf("UpdateUtilMetrics error: %v", err)
	}
}

var _ MetricPoller = (*Poller)(nil)
