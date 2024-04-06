package e

import "errors"

func WrapError(msg string, err error) error {
	return errors.New(msg + err.Error())
}
