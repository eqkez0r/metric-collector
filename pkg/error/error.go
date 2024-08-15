// Пакет для обертки ошибок
package e

import "errors"

// Функция WrapError оборачивает ошибку
func WrapError(msg string, err error) error {
	if err == nil {
		return nil
	}
	return errors.New(msg + err.Error())
}
