package poster

import (
	"github.com/Eqke/metric-collector/internal/agent/reqtype"
	"github.com/Eqke/metric-collector/internal/agent/result"
	"github.com/Eqke/metric-collector/internal/config"
	"go.uber.org/zap"
	"sync"
)

type poster struct {
	settings *config.AgentConfig
	logger   *zap.SugaredLogger
	errChan  chan error
	res      *result.Result
}

func NewPoster(
	logger *zap.SugaredLogger,
	settings *config.AgentConfig,
) *poster {
	return &poster{
		logger:   logger,
		settings: settings,
		errChan:  make(chan error),
		res:      result.New(),
	}
}

func (p *poster) Shutdown() {
	close(p.errChan)
}

func (p *poster) Post(requests <-chan *reqtype.ReqType) {

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

func (p *poster) errorLogger(done chan struct{}) {
	for {
		select {
		case <-done:
			return
		case err := <-p.errChan:
			p.logger.Error(err)
		default:
		}

	}
}

func (p *poster) postRequest(r *reqtype.ReqType) {

	resp, err := r.Req.Post(r.Endpoint)
	p.res.IncAll()
	if err != nil {
		p.res.IncErrors()
		p.errChan <- err
	}
	p.logger.Infof("Request status: %s", resp.Status())
}
