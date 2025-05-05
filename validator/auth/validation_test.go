package auth

import (
	"os"
	"testing"

	authDto "github.com/mohamadrezamomeni/momo/dto/controller/auth"
)

var validator *Validation

func TestMain(m *testing.M) {
	validator = New()
	code := m.Run()
	os.Exit(code)
}

func TestLogin(t *testing.T) {
	err := validator.LoginValidator(authDto.Login{
		Username: "mohamadreza",
		Password: "Mamad@1234",
	})
	if err != nil {
		t.Error("error has occured")
	}

	err = validator.LoginValidator(authDto.Login{
		Username: "",
		Password: "Mamad@1234",
	})

	if err == nil {
		t.Error("validator doesn't work very well")
	}
}
