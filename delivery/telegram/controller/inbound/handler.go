package inbound

import (
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
	inboundValidator "github.com/mohamadrezamomeni/momo/validator/inbound"
)

type Handler struct {
	userSvc          UserService
	inboundSvc       InboundService
	inboundValidator *inboundValidator.Validator
}

type InboundService interface {
	Filter(inpt *inboundServiceDto.FilterInbounds) ([]*entity.Inbound, error)
	Create(inpt *inboundServiceDto.CreateInbound) (*entity.Inbound, error)
	FindInboundByID(string) (*entity.Inbound, error)
}

type UserService interface {
	FindByTelegramID(string) (*entity.User, error)
}

func New(
	userSvc UserService,
	inboundSvc InboundService,
	inboundValidator *inboundValidator.Validator,
) *Handler {
	return &Handler{
		userSvc:          userSvc,
		inboundSvc:       inboundSvc,
		inboundValidator: inboundValidator,
	}
}
