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

func Fatalf(format string, args ...any) {
	momoLogger.Debuggingf(format, args...)

	os.Exit(1)
}

func Fatal(s string) {
	momoLogger.Debugging(s)

	os.Exit(1)
}

func Errorf(format string, args ...any) error {
	momoLogger.Warrningf(format, args...)
	s := fmt.Sprintf(format, args...)
	return &MomoError{
		args:    args,
		pattern: format,
		message: s,
	}
}

func Error(s string) error {
	momoLogger.Warrning(s)
	return &MomoError{
		pattern: s,
		message: s,
	}
}

func DebuggingErrorf(format string, args ...any) error {
	momoLogger.Debuggingf(format, args...)
	return &MomoError{
		args:    args,
		pattern: format,
		message: fmt.Sprintf(format, args...),
	}
}

func DebuggingError(s string) error {
	momoLogger.Debugging(s)
	return &MomoError{
		pattern: s,
		message: s,
	}
}
