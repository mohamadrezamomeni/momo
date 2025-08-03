package charge

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	chargeControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/charge"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (v *Validator) ValidateFilterCharges(req chargeControllerDto.FilterCharges) error {
	scope := "validation.ValidateFilterCharges"

	err := validation.ValidateStruct(
		&req,
		validation.Field(
			&req.UserID,
			validation.When(req.UserID != "", is.UUID),
		),
		validation.Field(
			&req.InboundID,
			validation.When(req.UserID != "", is.Int),
		),
		validation.Field(
			&req.Statuses,
			validation.When(req.Statuses != "", validation.By(func(value interface{}) error {
				labels := strings.Split(value.(string), ",")
				return validation.Validate(labels,
					validation.Each(
						validation.In(
							"approved",
							"pending",
							"rejected",
							"assigned",
							"unknown",
						),
					),
				)
			})),
		),
	)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).BadRequest().Input(req).ErrorWrite()
	}
	return nil
}
