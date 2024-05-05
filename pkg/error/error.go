package e

import "errors"

func WrapError(msg string, err error) error {
	if err == nil {
		return nil
	}
	return errors.New(msg + err.Error())
}
