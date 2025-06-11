package vpnsource

import (
	VPNSourceServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn_source"
	"github.com/mohamadrezamomeni/momo/service/auth"
)

type Handler struct {
	VPNSourceSvc VPNSourceService
	authSvc      *auth.Auth
}

type VPNSourceService interface {
	Create(*VPNSourceServiceDto.CreateVPNSourceDto) error
}

func New(VPNSourceService VPNSourceService, authSvc *auth.Auth) *Handler {
	return &Handler{
		VPNSourceSvc: VPNSourceService,
		authSvc:      authSvc,
	}
}
