package httperror

import (
	"fmt"
	"net/http"
	"testing"

	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func TestStatus(t *testing.T) {
	code := getStatus(momoError.Scope("test"))
	if code != http.StatusInternalServerError {
		t.Error("code must be 500")
	}

	code = getStatus(momoError.Scope("test").Forbidden())
	if code != http.StatusForbidden {
		t.Error("code must be 403")
	}

	code = getStatus(momoError.Scope("test").BadRequest())
	if code != http.StatusBadRequest {
		t.Error("code must be 400")
	}

	code = getStatus(fmt.Errorf("error"))
	if code != http.StatusInternalServerError {
		t.Error("code must be 500")
	}
}

func TestMessage(t *testing.T) {
	message := getMessage(fmt.Errorf(""))
	if message != "something went wrong" {
		t.Error("error is incompatible")
	}

	message = getMessage(momoError.Scope(""))
	if message != "something went wrong" {
		t.Error("error is incompatible")
	}

	msg1 := "hello world"
	message = getMessage(momoError.Scope("").Errorf(msg1))
	if message != msg1 {
		t.Error("error is incompatible")
	}
}
