package error

import "errors"

func New(errString string) error {
	return errors.New(errString)
}
