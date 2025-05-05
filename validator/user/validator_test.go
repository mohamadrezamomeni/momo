package user

import (
	"os"
	"testing"

	"github.com/mohamadrezamomeni/momo/dto/controller/user"
)

var userValidator *Validator

func TestMain(m *testing.M) {
	userValidator = New()
	code := m.Run()
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	err := userValidator.ValidateAddUserRequest(user.AddUser{
		IsAdmin:   true,
		Username:  "jackson",
		FirstName: "micheal",
		LastName:  "jackson",
		Password:  "Mamad@1234",
	})
	if err != nil {
		t.Errorf("some thing wrong has happend the problem was %v", err)
	}

	err = userValidator.ValidateAddUserRequest(user.AddUser{
		IsAdmin:   true,
		Username:  "jackson",
		FirstName: "mich",
		LastName:  "jackson",
		Password:  "mamaaaaa",
	})
	if err == nil {
		t.Errorf("err Expectect password validation works very well ")
	}
}
