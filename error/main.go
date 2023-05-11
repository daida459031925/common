package error

import (
	"errors"
	"github.com/daida459031925/common/fmt"
	"log"
)

func New(errString string) error {
	return errors.New(errString)
}

func NewSprintf(errString string, args ...any) error {
	return errors.New(fmt.Sprintf(errString, args...))
}

func RuntimeException(e error) {
	panic(e)
}

func RuntimeExceptionTF(tf bool, e error) {
	if tf {
		panic(e)
	} else {
		log.Fatal(e)
	}
}
