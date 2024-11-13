package restorer

import (
	"context"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"
)

type ToJSONProvider interface {
	ToJSON(context.Context) ([]byte, error)
}

type Restorer struct {
	logger   *zap.SugaredLogger
	storage  ToJSONProvider
	ticker   time.Duration
	filepath string
}

func New(
	logger *zap.SugaredLogger,
	storage ToJSONProvider,
	filepath string,
	duration int,
) *Restorer {
	return &Restorer{
		logger:   logger,
		storage:  storage,
		ticker:   time.Duration(duration) * time.Second,
		filepath: filepath,
	}
}

func (r *Restorer) Run(
	ctx context.Context,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	r.logger.Info("Restore was started")
	t := time.NewTicker(r.ticker)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			{
				r.logger.Info("Restore was finished")
				return
			}
		case <-t.C:
			{
				r.logger.Info("Restored...")
				r.restore(ctx)
				r.logger.Info("Restore was finished")
			}
		}
	}
}

func (r *Restorer) restore(ctx context.Context) {
	func() {
		b, err := r.storage.ToJSON(ctx)
		if err != nil {
			r.logger.Errorf("Restore error getting storage: %v", err)
			return
		}
		f, err := os.OpenFile(r.filepath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			r.logger.Errorf("Restore error opening file: %v", err)
			return
		}
		defer f.Close()
		_, err = f.Write(b)
		if err != nil {
			r.logger.Errorf("Restore error writing file: %v", err)
			return
		}
	}()
}
