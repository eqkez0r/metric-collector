package retry

import "time"

func Retry(attempts int, f func() error) error {
	for i := 0; i < attempts; i++ {
		if err := f(); err == nil {
			return nil
		}
		time.Sleep(time.Second * time.Duration(1+2*attempts))
	}
	return nil
}
