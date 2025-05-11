package inbound

import (
	"time"

	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
	authSvc "github.com/mohamadrezamomeni/momo/service/auth"
	inboundValidation "github.com/mohamadrezamomeni/momo/validator/inbound"
)

type Handler struct {
	inboundSvc InboundService
	validation *inboundValidation.Validator
	authSvc    *authSvc.Auth
}

type InboundService interface {
	Create(*inboundServiceDto.CreateInbound) (*entity.Inbound, error)
	Filter(*inboundServiceDto.FilterInbounds) ([]*entity.Inbound, error)
	Block(string) error
	UnBlock(string) error
	ExtendInbound(string, time.Time) error
}

func New(inboundSvc InboundService, validation *inboundValidation.Validator, authSvc *authSvc.Auth) *Handler {
	return &Handler{
		inboundSvc: inboundSvc,
		validation: validation,
		authSvc:    authSvc,
	}
}
