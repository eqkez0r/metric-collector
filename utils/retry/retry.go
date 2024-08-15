// Пакет retry содержит функцию, которая выполняет переданную функцию n-раз
// до тех пор, пока ошибка не будет равна nil
package retry

import (
	"time"

	"go.uber.org/zap"
)

// Функция Retry получает кол-во повторов и необходимую функцию
func Retry(
	logger *zap.SugaredLogger,
	attempts int,
	f func() error) error {
	// Выполнение функции
	if err := f(); err != nil {
		// Если получена ошибка, то начинаем механизм повторов
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
