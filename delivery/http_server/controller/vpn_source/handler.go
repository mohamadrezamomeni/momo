package vpnsource

import (
	VPNSourceServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn_source"
	"github.com/mohamadrezamomeni/momo/entity"
	"github.com/mohamadrezamomeni/momo/service/auth"
	vpnSourceValidation "github.com/mohamadrezamomeni/momo/validator/vpn_source"
)

type Handler struct {
	VPNSourceSvc       VPNSourceService
	authSvc            *auth.Auth
	vpnSourceValidator *vpnSourceValidation.Validator
}

type VPNSourceService interface {
	Create(*VPNSourceServiceDto.CreateVPNSourceDto) error
	FilterVPNSources(*VPNSourceServiceDto.FilterVPNSourcesDto) ([]*entity.VPNSource, error)
}

func New(
	VPNSourceService VPNSourceService,
	vpnSourceValidator *vpnSourceValidation.Validator,
	authSvc *auth.Auth,
) *Handler {
	return &Handler{
		VPNSourceSvc:       VPNSourceService,
		authSvc:            authSvc,
		vpnSourceValidator: vpnSourceValidator,
	}
}
