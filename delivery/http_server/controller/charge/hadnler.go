package charge

import (
	chargeServiceDto "github.com/mohamadrezamomeni/momo/dto/service/charge"
	"github.com/mohamadrezamomeni/momo/entity"
	"github.com/mohamadrezamomeni/momo/service/auth"
	validator "github.com/mohamadrezamomeni/momo/validator/charge"
)

type Handler struct {
	chargeSvc  ChargeService
	authSvc    *auth.Auth
	validation *validator.Validator
}

type ChargeService interface {
	ApproveCharge(string) error
	FilterCharges(*chargeServiceDto.FilterCharges) ([]*entity.Charge, error)
}

func New(
	chargeSvc ChargeService,
	authSvc *auth.Auth,
	validation *validator.Validator,
) *Handler {
	return &Handler{
		chargeSvc:  chargeSvc,
		authSvc:    authSvc,
		validation: validation,
	}
}
