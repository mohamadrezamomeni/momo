package vpnsource

import (
	VPNSourceServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn_source"
	"github.com/mohamadrezamomeni/momo/entity"
)

type Handler struct {
	vpnSourceService VPNSourceService
	userSvc          UserService
}

type UserService interface {
	FindByTelegramID(string) (*entity.User, error)
}

type VPNSourceService interface {
	FilterVPNSources(*VPNSourceServiceDto.FilterVPNSourcesDto) ([]*entity.VPNSource, error)
	Find(string) (*entity.VPNSource, error)
}

func New(vpnSourceService VPNSourceService, userSvc UserService) *Handler {
	return &Handler{
		vpnSourceService: vpnSourceService,
		userSvc:          userSvc,
	}
}
