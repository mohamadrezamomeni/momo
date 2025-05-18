package inbound

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	inboundControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/inbound"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	timeTransofrmer "github.com/mohamadrezamomeni/momo/transformer/time"
)

func (v *Validator) ValidateExtendingInbound(req inboundControllerDto.ExtendInboundDto) error {
	scope := "validator.ValidateExtendingInbound"
	err := validation.ValidateStruct(
		&req,
		validation.Field(
			&req.ExtendedTrafficLimit,
			validation.Required,
		),
		validation.Field(
			&req.End,
			validation.Required,
			validation.By(func(value interface{}) error {
				_, ok := value.(string)
				if !ok {
					return momoError.Scope(scope).Input(req).BadRequest().ErrorWrite()
				}
				return nil
			}),
		),
	)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).BadRequest().Input(req).ErrorWrite()
	}

	inbound, err := v.inboundSvc.FindInboundByID(req.ID)
	if err != nil {
		return err
	}

	endTime, err := timeTransofrmer.ConvertStrToTime(req.End)
	if err != nil {
		return err
	}

	now := time.Now()
	if inbound.End.After(endTime) || now.After(inbound.End) {
		return momoError.Scope(scope).Input(req).BadRequest().ErrorWrite()
	}

	if inbound.TrafficLimit > inbound.TrafficUsage {
		return momoError.Scope(scope).Input(req).BadRequest().ErrorWrite()
	}
	return nil
}
