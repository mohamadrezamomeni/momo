package vpnsource

import (
	VPNSourceServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn_source"
	"github.com/mohamadrezamomeni/momo/entity"
)

type Handler struct {
	vpnSourceService VPNSourceService
}

type VPNSourceService interface {
	FilterVPNSources(*VPNSourceServiceDto.FilterVPNSourcesDto) ([]*entity.VPNSource, error)
	Find(string) (*entity.VPNSource, error)
}

func New(vpnSourceService VPNSourceService) *Handler {
	return &Handler{
		vpnSourceService: vpnSourceService,
	}
}
