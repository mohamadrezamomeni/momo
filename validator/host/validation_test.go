package host

import (
	"os"
	"testing"

	hostDto "github.com/mohamadrezamomeni/momo/dto/controller/host"
)

var validator *Validator

func TestMain(m *testing.M) {
	validator = New()
	code := m.Run()
	os.Exit(code)
}

func TestCreateHostValidation(t *testing.T) {
	err := validator.CreateHostValidation(hostDto.CreateHostDto{
		Domain: "twitter.com",
		Port:   "12345",
	})
	if err != nil {
		t.Fatalf("the problem has occurred")
	}
	err = validator.CreateHostValidation(hostDto.CreateHostDto{
		Domain: "twitter.com",
		Port:   "1234f5",
	})
	if err == nil {
		t.Fatal("the creating host validation has occured")
	}
	err = validator.CreateHostValidation(hostDto.CreateHostDto{
		Port: "1234f5",
	})
	if err == nil {
		t.Fatal("the creating host validation has occured")
	}
}
