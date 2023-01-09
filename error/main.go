package error

import "errors"

func New(errString string) error {
	return errors.New(errString)
}

func RuntimeException(e error) {
	panic(e)
}

func RuntimeExceptionTF(tf bool, e error) {
	if tf {
		panic(e)
	}
}
