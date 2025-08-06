package charge

import (
	chargeServiceDto "github.com/mohamadrezamomeni/momo/dto/service/charge"
	"github.com/mohamadrezamomeni/momo/entity"
	"github.com/mohamadrezamomeni/momo/service/auth"
	validator "github.com/mohamadrezamomeni/momo/validator/charge"
)

type Handler struct {
	chargeSvc        ChargeService
	authSvc          *auth.Auth
	validation       *validator.Validator
	inboundChargeSvc InboundChargeService
}

type ChargeService interface {
	FilterCharges(*chargeServiceDto.FilterCharges) ([]*entity.Charge, error)
}

type InboundChargeService interface {
	ApproveCharge(string) error
}

func New(
	chargeSvc ChargeService,
	authSvc *auth.Auth,
	validation *validator.Validator,
	inboundChargeSvc InboundChargeService,
) *Handler {
	return &Handler{
		inboundChargeSvc: inboundChargeSvc,
		chargeSvc:        chargeSvc,
		authSvc:          authSvc,
		validation:       validation,
	}
}
