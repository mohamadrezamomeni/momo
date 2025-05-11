package inbound

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	inboundControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/inbound"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (v *Validator) ValidateExtendingInbound(req inboundControllerDto.ExtendInboundDto) error {
	scope := "validator.ValidateExtendingInbound"
	err := validation.ValidateStruct(
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

	endTime := utils.GetDateTime(req.End)
	now := time.Now()
	fmt.Println(inbound.End, endTime, inbound.End.After(endTime))
	if inbound.End.After(endTime) || now.After(inbound.End) {
		return momoError.Scope(scope).BadRequest().ErrorWrite()
	}
	return nil
}
