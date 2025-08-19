package charge

import (
	chargeServiceDto "github.com/mohamadrezamomeni/momo/dto/service/charge"
	"github.com/mohamadrezamomeni/momo/entity"
)

type Handler struct {
	chargeSvc ChargeService
}

type ChargeService interface {
	Create(*chargeServiceDto.CreateChargeDto) (*entity.Charge, error)
	FilterCharges(*chargeServiceDto.FilterCharges) ([]*entity.Charge, error)
}

func New(chargeSvc ChargeService) *Handler {
	return &Handler{
		chargeSvc: chargeSvc,
	}
}
