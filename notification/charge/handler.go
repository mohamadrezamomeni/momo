package charge

import "github.com/mohamadrezamomeni/momo/entity"

type Handler struct {
	userSvc    UserService
	inboundSvc InboundService
	chargeSvc  ChargeService
}

type UserService interface {
	FindByID(id string) (*entity.User, error)
}

type InboundService interface {
	FindInboundByID(string) (*entity.Inbound, error)
}

type ChargeService interface {
	FindByID(string) (*entity.Charge, error)
}

func New(userSvc UserService, inboundSvc InboundService, chargeSvc ChargeService) *Handler {
	return &Handler{
		userSvc:    userSvc,
		chargeSvc:  chargeSvc,
		inboundSvc: inboundSvc,
	}
}
