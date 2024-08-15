// Пакет poller отвечает за сбор метрик
package poller

import (
	"runtime"
	"sync"

	"github.com/Eqke/metric-collector/pkg/metric"
	"go.uber.org/zap"
)

// Интерфейс MetricPoller отвечает за получение метрик
type MetricPoller interface {
	Poll() metric.Map
}

// Тип Poller является реализацией интерфейса MetricPoller
type Poller struct {
	logger *zap.SugaredLogger
	mp     map[metric.MType]map[metric.Name]string
	mu     sync.Mutex
	wg     sync.WaitGroup
}

// Функция NewPoller возвращает объект типа Poller
func NewPoller(logger *zap.SugaredLogger) *Poller {
	mp := make(metric.Map)
	mp[metric.TypeGauge] = make(map[metric.Name]string)
	mp[metric.TypeCounter] = make(map[metric.Name]string)
	return &Poller{
		logger: logger,
		mp:     mp,
		mu:     sync.Mutex{},
		wg:     sync.WaitGroup{},
	}
}

// Метод Poll собирает метрики и возвращает объект metric.Map
func (p *Poller) Poll() metric.Map {

	p.wg.Add(2)
	go p.updateRuntime()
	go p.updateUtil()
	p.wg.Wait()

	return p.mp
}

// Метод updateRuntime обновляет метрики из пакета runtime
func (p *Poller) updateRuntime() {
	defer p.wg.Done()
	ms := &runtime.MemStats{}
	runtime.ReadMemStats(ms)
	metric.UpdateRuntimeMetrics(ms, p.mp)

}

// Метод updateUtil обновляет метрики cpu, totalmemory и freememory
func (p *Poller) updateUtil() {
	defer p.wg.Done()
	if err := metric.UpdateUtilMetrics(p.mp); err != nil {
		p.logger.Errorf("UpdateUtilMetrics error: %v", err)
	}
}

var _ MetricPoller = (*Poller)(nil)
