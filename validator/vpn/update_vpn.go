package vpn

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	vpnControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/vpn"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (v *Validator) ValidateUpdatingVPN(req vpnControllerDto.UpdateVPN) error {
	scope := "validation.Filter"
	err := validation.ValidateStruct(
		&req,
		validation.Field(
			&req.VPNStatusLabel,
			validation.When(req.VPNStatusLabel != "", validation.In(
				"cordon",
				"drain",
				"ready",
				"unkhown",
			)),
		),
	)
	if err != nil {
		return momoError.Wrap(err).BadRequest().Scope(scope).Input(req).ErrorWrite()
	}
	return nil
}
