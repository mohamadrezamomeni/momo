package user

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	userControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/user"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (v *Validator) ValidateAddUserRequest(data *userControllerDto.AddUser) error {
	scope := "validator.user.ValidateAddUserRequest"

	err := validation.ValidateStruct(data,
		validation.Field(
			data.Password,
			validation.Required,
			validation.Length(6, 15),
			validation.Match(regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*])[A-Za-z\d!@#$%^&*]{8,}$`)).
				Error("must be at least 8 characters and include upper, lower, digit, and special character"),
		),
		validation.Field(
			data.FirstName,
			validation.Required,
			validation.Length(2, 20),
		),
		validation.Field(
			data.LastName,
			validation.Required,
			validation.Length(2, 20),
		),
		validation.Field(
			data.IsAdmin,
			validation.Required,
			validation.In(true, false),
		),
	)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Errorf("the input is %+v", *data)
	}
	return err
}
