package host

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	hostDto "github.com/mohamadrezamomeni/momo/dto/controller/host"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (v *Validator) CreateHostValidation(req hostDto.CreateHostDto) error {
	scope := "validation.CreateHostValidation"

	err := validation.ValidateStruct(
		&req,
		validation.Field(
			&req.Domain,
			validation.Required,
			validation.Length(3, 0),
		),
		validation.Field(
			&req.Port,
			validation.Required,
			validation.Match(regexp.MustCompile(`^\d+$`)).Error("must be a numeric string"),
		),
	)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).BadRequest().Input(req).ErrorWrite()
	}
	return nil
}
