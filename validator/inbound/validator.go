package inbound

import "github.com/mohamadrezamomeni/momo/entity"

type Validator struct {
	userSvc    UserService
	inboundSvc InboundService
}

type UserService interface {
	FindByID(string) (*entity.User, error)
}

type InboundService interface {
	FindInboundByID(string) (*entity.Inbound, error)
}

func New(userSvc UserService, inboundSvc InboundService) *Validator {
	return &Validator{
		userSvc:    userSvc,
		inboundSvc: inboundSvc,
	}
}
