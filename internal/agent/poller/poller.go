package poller

import (
	"github.com/Eqke/metric-collector/pkg/metric"
	"go.uber.org/zap"
	"runtime"
	"sync"
)

type poller struct {
	logger *zap.SugaredLogger
	mp     map[metric.TypeMetric]map[metric.MetricName]string
	mu     sync.Mutex
	wg     sync.WaitGroup
}

func NewPoller(logger *zap.SugaredLogger) *poller {
	mp := make(metric.MetricMap)
	mp[metric.TypeGauge] = make(map[metric.MetricName]string)
	mp[metric.TypeCounter] = make(map[metric.MetricName]string)
	return &poller{
		logger: logger,
		mp:     mp,
		mu:     sync.Mutex{},
		wg:     sync.WaitGroup{},
	}
}

func (p *poller) Poll() metric.MetricMap {

	p.wg.Add(2)
	go p.updateRuntime()
	go p.updateUtil()
	p.wg.Wait()

	return p.mp
}

func (p *poller) updateRuntime() {
	defer p.wg.Done()
	ms := &runtime.MemStats{}
	runtime.ReadMemStats(ms)
	metric.UpdateRuntimeMetrics(ms, p.mp)

}

func (p *poller) updateUtil() {
	defer p.wg.Done()
	if err := metric.UpdateUtilMetrics(p.mp); err != nil {
		p.logger.Errorf("UpdateUtilMetrics error: %v", err)
	}
}
