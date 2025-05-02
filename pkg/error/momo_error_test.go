package error

import (
	"fmt"
	"testing"
)

func TestWithoutMainError(t *testing.T) {
	scopeTest := "test.TestWithoutMainError"
	err := Scope(scopeTest).DeactiveWrite().DebuggingErrorf("patern was with arguments %d", 1)

	message := "the scope is \"test.TestWithoutMainError\" and the main error is \"nothing\" the additional information is \"patern was with arguments 1\""

	if err.Error() != message {
		t.Error("error to compare we expected error and the error we were given")
	}

	err = Scope(scopeTest).DeactiveWrite().DebuggingErrorf("patern was without any arguments")

	message = "the scope is \"test.TestWithoutMainError\" and the main error is \"nothing\" the additional information is \"patern was without any arguments\""
	if err.Error() != message {
		t.Error("error to compare we expected error and the error we were given")
	}

	err = Wrap(fmt.Errorf("database error")).DeactiveWrite().DebuggingErrorf("patern was without any arguments")
	message = "the scope is \"empty\" and the main error is \"database error\" the additional information is \"patern was without any arguments\""
	if err.Error() != message {
		t.Error("error to compare we expected error and the error we were given")
	}

	err = Wrap(fmt.Errorf("database error")).Scope(scopeTest).DeactiveWrite().DebuggingErrorf("patern was without any arguments")
	message = "the scope is \"test.TestWithoutMainError\" and the main error is \"database error\" the additional information is \"patern was without any arguments\""
	if err.Error() != message {
		t.Error("error to compare we expected error and the error we were given")
	}
}
