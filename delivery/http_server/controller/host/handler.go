package host

import (
	hostServiceDto "github.com/mohamadrezamomeni/momo/dto/service/host"
	"github.com/mohamadrezamomeni/momo/entity"
	authService "github.com/mohamadrezamomeni/momo/service/auth"
	hostValidation "github.com/mohamadrezamomeni/momo/validator/host"
)

type Handler struct {
	hostSvc       HostService
	authSvc       *authService.Auth
	hostValidator *hostValidation.Validator
}

type HostService interface {
	Create(*hostServiceDto.CreateHostDto) (*entity.Host, error)
	Filter(*hostServiceDto.FilterHosts) ([]*entity.Host, error)
}

func New(hostSvc HostService, authSvc *authService.Auth, hostValidation *hostValidation.Validator) *Handler {
	return &Handler{
		hostValidator: hostValidation,
		hostSvc:       hostSvc,
		authSvc:       authSvc,
	}
}
