package inbound

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	inboundControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/inbound"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
	timetransfer "github.com/mohamadrezamomeni/momo/transformer/time"
)

func (v *Validator) ValidateSettingPeriodTime(req inboundControllerDto.SetPeriodDto) error {
	scope := "validator.ValidateSettingPeriodTime"
	err := validation.ValidateStruct(
		&req,
		validation.Field(
			&req.TrafficLimit,
			validation.Required,
		),
		validation.Field(
			&req.Start,
			validation.Required,
			validation.By(func(value interface{}) error {
				s, ok := value.(string)
				if !ok {
					return momoError.Scope(scope).Input(req).BadRequest().ErrorWrite()
				}
				startTime, err := timetransfer.ConvertStrToTime(s)
				if err != nil {
					return err
				}

				if !ok {
					return momoError.Scope(scope).Input(req).BadRequest().ErrorWrite()
				}
				endTime, err := timetransfer.ConvertStrToTime(req.End)
				if err != nil {
					return nil
				}
				if startTime.After(endTime) {
					return momoError.Scope(scope).Errorf("start time must be lower than end time")
				}
				return nil
			}),
		),
		validation.Field(
			&req.End,
			validation.Required,
		),
	)
	if err != nil {
		return momoError.Wrap(err).Scope(scope).BadRequest().Input(req).ErrorWrite()
	}

	inbound, err := v.inboundSvc.FindInboundByID(req.ID)
	if err != nil {
		return err
	}
	now := time.Now()
	startTime := utils.GetDateTime(req.Start)

	if !now.After(inbound.End) || now.After(startTime) {
		return momoError.Scope(scope).Input(req).ErrorWrite()
	}

	return nil
}
