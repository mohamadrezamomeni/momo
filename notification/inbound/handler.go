package inbound

import "github.com/mohamadrezamomeni/momo/entity"

type Handler struct {
	userSvc    UserService
	inboundSvc InboundService
}

type UserService interface {
	FindByID(string) (*entity.User, error)
}

type InboundService interface {
	FindInboundByID(string) (*entity.Inbound, error)
}

func New(userSvc UserService, inboundSvc InboundService) *Handler {
	return &Handler{
		userSvc:    userSvc,
		inboundSvc: inboundSvc,
	}
}
