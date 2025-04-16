package error

import (
	"fmt"
	"os"

	momoLogger "momo/pkg/log"
)

type MomoError struct {
	args    []interface{}
	pattern string
	message string
}

func (m *MomoError) Error() string {
	return m.message
}

func Fatalf(format string, args ...interface{}) {
	momoLogger.Warrning(fmt.Sprintf(format, args))

	os.Exit(1)
}

func Errorf(format string, args ...interface{}) error {
	s := fmt.Sprintf(format, args)
	momoLogger.Warrning(s)
	return &MomoError{
		args:    args,
		pattern: format,
		message: s,
	}
}

func Error(format string) error {
	momoLogger.Warrning(format)
	return &MomoError{
		pattern: format,
		message: format,
	}
}
