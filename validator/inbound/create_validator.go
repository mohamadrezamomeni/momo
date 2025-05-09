package inbound

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	inboundControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	timetransfer "github.com/mohamadrezamomeni/momo/transformer/time"
)

func (v *Validator) ValidateCreatingInbound(req inboundControllerDto.CreateInbound) error {
	scope := "validator.ValidateCreatingInbound"
	err := validation.ValidateStruct(
		&req,
		validation.Field(
			&req.Domain,
			validation.Required,
			is.URL,
		),
		validation.Field(
			&req.Port,
			validation.Required,
			validation.Match(regexp.MustCompile(`^\d+$`)).Error("must be a numeric string"),
		),
		validation.Field(
			&req.UserID,
			validation.Required,
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
					return momoError.Scope(scope).Input(req).ErrorWrite()
				}
				if entity.UknownVPNType == entity.ConvertStringVPNTypeToEnum(v) {
					return momoError.Scope(scope).Input(req).BadRequest().ErrorWrite()
				}
				return nil
			}),
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
	user, err := v.userSvc.FindByID(req.UserID)
	if err != nil || user == nil {
		return momoError.Wrap(err).Scope(scope).Input(req).BadRequest().Errorf("error to retrive user")
	}
	return nil
}
