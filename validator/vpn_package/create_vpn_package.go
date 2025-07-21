package vpnpackage

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	vpnPackageDto "github.com/mohamadrezamomeni/momo/dto/controller/vpn_package"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (v *Validator) CreateVPNPackage(req vpnPackageDto.CreateVPNPackage) error {
	scope := "vpnPackage.Create"
	err := validation.ValidateStruct(
		&req,
		validation.Field(
			&req.PriceTitle,
			validation.Required,
		),
		validation.Field(
			&req.TrafficLimitTitle,
			validation.Required,
		),
	)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).BadRequest().Input(req).ErrorWrite()
	}
	return nil
}
