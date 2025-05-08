package vpn

import (
	vpnServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn"
	"github.com/mohamadrezamomeni/momo/entity"
	"github.com/mohamadrezamomeni/momo/service/auth"
	validator "github.com/mohamadrezamomeni/momo/validator/vpn"
)

type Handler struct {
	vpnSvc     VPNService
	validation *validator.Validator
	authSvc    *auth.Auth
}

type VPNService interface {
	Create(*vpnServiceDto.CreateVPN) (*entity.VPN, error)
	Filter(*vpnServiceDto.FilterVPNs) ([]*entity.VPN, error)
}

func New(vpnSvc VPNService, validation *validator.Validator, authSvc *auth.Auth) *Handler {
	return &Handler{
		vpnSvc:     vpnSvc,
		validation: validation,
		authSvc:    authSvc,
	}
}
