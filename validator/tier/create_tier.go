package charge

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	tierControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/tier"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (v *Validator) ValidateCreatingTier(req tierControllerDto.CreateTier) error {
	scope := "validation.ValidateCreatingTier"

	err := validation.ValidateStruct(
		&req,
		validation.Field(
			&req.Name,
			validation.Length(3, 0),
		),
	)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).BadRequest().Input(req).ErrorWrite()
	}
	return nil
}
