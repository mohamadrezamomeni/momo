package charge

import (
	chargeServiceDto "github.com/mohamadrezamomeni/momo/dto/service/charge"
	"github.com/mohamadrezamomeni/momo/entity"
)

type Handler struct {
	chargeSvc ChargeService
	userSvc   UserService
}

type ChargeService interface {
	Create(*chargeServiceDto.CreateChargeDto) (*entity.Charge, error)
	FilterCharges(*chargeServiceDto.FilterCharges) ([]*entity.Charge, error)
}

type UserService interface {
	FindByTelegramID(string) (*entity.User, error)
}

func New(chargeSvc ChargeService, userSvc UserService) *Handler {
	return &Handler{
		chargeSvc: chargeSvc,
		userSvc:   userSvc,
	}
}
