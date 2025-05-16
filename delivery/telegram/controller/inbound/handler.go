package inbound

import (
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
)

type Handler struct {
	userSvc    UserService
	inboundSvc InboundService
}

type InboundService interface {
	Filter(inpt *inboundServiceDto.FilterInbounds) ([]*entity.Inbound, error)
}

type UserService interface {
	FindByTelegramID(string) (*entity.User, error)
}

func New(userSvc UserService, inboundSvc InboundService) *Handler {
	return &Handler{
		userSvc:    userSvc,
		inboundSvc: inboundSvc,
	}
}
