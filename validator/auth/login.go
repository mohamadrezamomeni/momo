package auth

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	authDto "github.com/mohamadrezamomeni/momo/dto/controller/auth"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (v *Validation) LoginValidator(req authDto.Login) error {
	scope := "valid ation.login"
	err := validation.ValidateStruct(
		&req,
		validation.Field(&req.Username, validation.Required, validation.Length(1, 0)),
		validation.Field(&req.Password, validation.Required, validation.Length(1, 0)),
	)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).Input(req).BadRequest().ErrorWrite()
	}
	return nil
}
