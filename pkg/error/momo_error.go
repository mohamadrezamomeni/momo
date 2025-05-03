package error

import (
	"fmt"
	"os"

	momoLogger "github.com/mohamadrezamomeni/momo/pkg/log"
)

const empty = "empty"

type MomoError struct {
	args      []any
	pattern   string
	scope     string
	err       error
	isPrinted bool
}

func Scope(scope string) *MomoError {
	return &MomoError{
		isPrinted: true,
		args:      []any{},
		pattern:   "",
		err:       nil,
		scope:     fmt.Sprintf("\"%s\"", scope),
	}
}

func Wrap(err error) *MomoError {
	return &MomoError{
		isPrinted: true,
		args:      []any{},
		pattern:   "",
		err:       err,
		scope:     fmt.Sprintf("\"%s\"", empty),
	}
}

func (m *MomoError) DeactiveWrite() *MomoError {
	m.isPrinted = false
	return m
}

func (m *MomoError) ActiveWrite() *MomoError {
	m.isPrinted = true
	return m
}

func (m *MomoError) Scope(scope string) *MomoError {
	m.scope = fmt.Sprintf("\"%s\"", scope)
	return m
}

func (m *MomoError) Error() string {
	message := fmt.Sprintf("the scope is %s and the main error is \"%s\"", m.scope, m.mainError())

	additionlMessage := ""

	if len(m.pattern) > 0 && len(m.args) > 0 {
		additionlMessage = fmt.Sprintf(m.pattern, m.args...)
	} else if len(m.pattern) > 0 {
		additionlMessage = m.pattern
	}

	if len(additionlMessage) > 0 {
		message += " the additional information is " + `"` + additionlMessage + `"`
	}
	return message
}

func (m *MomoError) mainError() string {
	err, ok := m.err.(*MomoError)

	if ok {
		return err.mainError()
	}

	if m.err != nil {
		return m.err.Error()
	}
	return "nothing"
}

func (m *MomoError) Fatalf(pattern string, args ...any) {
	m.args = args
	m.pattern = pattern
	if m.isPrinted {
		momoLogger.Debugging(m.Error())
	}
	os.Exit(1)
}

func (m *MomoError) Errorf(pattern string, args ...any) error {
	m.args = args
	m.pattern = pattern
	if m.isPrinted {
		momoLogger.Warrning(m.Error())
	}
	return m
}

func (m *MomoError) DebuggingErrorf(pattern string, args ...any) error {
	m.args = args
	m.pattern = pattern
	if m.isPrinted {
		momoLogger.Debugging(m.Error())
	}
	return m
}

func (m *MomoError) DebuggingError() *MomoError {
	if m.isPrinted {
		momoLogger.Debugging(m.Error())
	}
	return m
}

func (m *MomoError) Fatal() {
	if m.isPrinted {
		momoLogger.Debugging(m.Error())
	}
	os.Exit(1)
}

func (m *MomoError) ErrorWrite() error {
	if m.isPrinted {
		momoLogger.Warrning(m.Error())
	}
	return m
}
