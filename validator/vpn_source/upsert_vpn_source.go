package vpnsource

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	vpnSourceDto "github.com/mohamadrezamomeni/momo/dto/controller/vpn_source"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (v *Validator) ValidateUpsert(req vpnSourceDto.CreateVPNSourceDto) error {
	scope := "vpnSourceValidator.ValidateUpsert"

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Country,
			validation.Required,
			validation.Length(1, 0),
		),
		validation.Field(&req.English,
			validation.Required,
			validation.Length(1, 0),
		),
	)
	if err != nil {
		return momoError.Wrap(err).BadRequest().Scope(scope).BadRequest().Input(req).ErrorWrite()
	}
	return nil
}
