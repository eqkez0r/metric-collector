// package retry
package retry

import (
	"time"

	"go.uber.org/zap"
)

func Retry(
	logger *zap.SugaredLogger,
	attempts int,
	f func() error) error {
	if err := f(); err != nil {
		for i := 0; i < attempts; i++ {
			logger.Infof("attempt: %d", i+1)
			if err = f(); err == nil {
				return nil
			}
			time.Sleep(time.Second * time.Duration(1))
		}
		return err
	}
	return nil
}
