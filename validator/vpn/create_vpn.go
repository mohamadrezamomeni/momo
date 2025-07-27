package vpn

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	vpnControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/vpn"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (v *Validator) ValidateCreatingVPN(req vpnControllerDto.CreateVPN) error {
	scope := "validation.createVPN"

	err := validation.ValidateStruct(
		&req,
		validation.Field(
			&req.Domain,
			validation.Required,
			validation.Length(3, 0),
		),
		validation.Field(
			&req.Country,
			validation.Required,
			validation.Length(1, 0),
		),
		validation.Field(
			&req.Port,
			validation.Required,
			validation.Match(regexp.MustCompile(`^\d+$`)).Error("must be a numeric string"),
		),
		validation.Field(&req.EndPort, validation.Min(2000)),
		validation.Field(
			&req.StartPort, validation.By(func(value interface{}) error {
				startPort, _ := value.(int)
				if startPort > req.EndPort {
					return momoError.Scope(scope).Input(req).Errorf("end_port must be greather than start_port")
				}
				return nil
			})),
		validation.Field(
			&req.VpnType,
			validation.Required,
			validation.By(func(value interface{}) error {
				v, ok := value.(string)
				if !ok ||
					entity.UknownVPNType == entity.ConvertStringVPNTypeToEnum(v) {
					return momoError.Scope(scope).Input(req).Errorf("input of vpnType is wrong")
				}

				return nil
			}),
		),
	)
	if err != nil {
		return momoError.Wrap(err).BadRequest().Scope(scope).Input(req).ErrorWrite()
	}

	_, err = v.vpnSourceSvc.Find(req.Country)
	if err != nil {
		return err
	}
	return nil
}
