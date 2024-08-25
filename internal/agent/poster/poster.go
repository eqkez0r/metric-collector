// Пакет poster предоставляет функционал по публикации
// метрик на сервер
package poster

import (
	"github.com/Eqke/metric-collector/internal/agent/config"
	"sync"

	"github.com/Eqke/metric-collector/internal/agent/reqtype"
	"github.com/Eqke/metric-collector/internal/agent/result"
	"go.uber.org/zap"
)

// Интерфейс MetricPoster отвечает за публикацию метрик
// на сервер
type MetricPoster interface {
	Post(requests <-chan *reqtype.ReqType)
}

// Тип Poster является реализацией MetricPoster
type Poster struct {
	settings *config.AgentConfig
	logger   *zap.SugaredLogger
	errChan  chan error
	res      *result.Result
}

// Функция NewPoster инциализирует и возвращает объект Poster
func NewPoster(
	logger *zap.SugaredLogger,
	settings *config.AgentConfig,
) *Poster {
	return &Poster{
		logger:   logger,
		settings: settings,
		errChan:  make(chan error),
		res:      result.New(),
	}
}

// Метод Post отправляет запросы, полученные из канала.
func (p *Poster) Post(requests <-chan *reqtype.ReqType) {

	var wg sync.WaitGroup
	done := make(chan struct{})

	go p.errorLogger(done)
	for i := 0; i < p.settings.RateLimit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.postRequest(<-requests)
		}()
	}

	wg.Wait()
	p.logger.Infof("All requests were sent: %d. Errors: %d", p.res.All(), p.res.Errors())
	defer p.res.Reset()
	defer close(done)
}

// Метод Shutdown позволяет корректно завершить работу Poster
func (p *Poster) Shutdown() {
	close(p.errChan)
}

// Метод errorLogger отвечает за логгирование ошибок во время работы Poster
func (p *Poster) errorLogger(done chan struct{}) {
	for {
		select {
		case <-done:
			return
		case err := <-p.errChan:
			p.logger.Error(err)
		default:
			{
			}
		}

	}
}

// Метод postRequest отвечает за публикацию метрики на сервер
func (p *Poster) postRequest(r *reqtype.ReqType) {

	resp, err := r.Req.Post(r.Endpoint)
	p.res.IncAll()
	if err != nil {
		p.res.IncErrors()
		p.errChan <- err
	}
	p.logger.Infof("Request status: %s", resp.Status())
}

var _ MetricPoster = (*Poster)(nil)