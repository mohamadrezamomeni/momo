package inbound

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	inboundControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (v *Validator) ValidateFilteringInbounds(req inboundControllerDto.FilterInboundsDto) error {
	scope := "validation.ValidateFilteringInbounds"
	err := validation.ValidateStruct(
		&req,
		validation.Field(
			&req.Domain,
			is.URL,
		),
		validation.Field(
			&req.Port,
			validation.Match(regexp.MustCompile(`^\d+$`)).Error("must be a numeric string"),
		),
		validation.Field(
			&req.UserID,
			validation.Match(regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[a-f0-9]{12}$`)),
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
				if entity.UknownVPNType == entity.ConvertStringVPNTypeToEnum(v) || 0 == entity.ConvertStringVPNTypeToEnum(v) {
					return momoError.Scope(scope).BadRequest().Input(req).ErrorWrite()
				}
				return nil
			}),
		),
	)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).BadRequest().Input(req).ErrorWrite()
	}
	return nil
}
