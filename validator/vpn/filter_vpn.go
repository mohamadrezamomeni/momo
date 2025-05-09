package vpn

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	vpnControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/vpn"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (v *Validator) ValidateFilterVPNs(req vpnControllerDto.FilterVPNs) error {
	scope := "validation.Filter"
	err := validation.ValidateStruct(
		&req,
		validation.Field(
			&req.Domain,
			validation.Length(3, 0),
		),
		validation.Field(
			&req.VPNType,
			validation.By(func(value interface{}) error {
				if validation.IsEmpty(value) {
					return nil
				}

				v, ok := value.(string)
				if !ok {
					return momoError.Scope(scope).BadRequest().Input(req).ErrorWrite()
				}
				if entity.UknownVPNType == entity.ConvertStringVPNTypeToEnum(v) {
					return momoError.Scope(scope).BadRequest().Input(req).ErrorWrite()
				}
				return nil
			}),
		),
	)
	if err != nil {
		return momoError.Wrap(err).BadRequest().Scope(scope).Input(req).ErrorWrite()
	}
	return nil
}
