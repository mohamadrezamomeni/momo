package inbound

import (
	"github.com/mohamadrezamomeni/momo/adapter"
	inboundRepoDto "github.com/mohamadrezamomeni/momo/dto/repository/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
)

type InboundTraffic struct {
	inboundRepo InboundTrafficRepository
	vpnService  VpnTrafficService
	userService UserTrafficService
}

type InboundTrafficRepository interface {
	Filter(*inboundRepoDto.FilterInbound) ([]*entity.Inbound, error)
	IncreaseTrafficUsage(string, uint32) error
}

type VpnTrafficService interface {
	MakeProxy() (adapter.ProxyVPN, error)
}

type UserTrafficService interface {
	FindByID(string) (*entity.User, error)
}

func NewInboundTraffic(
	inboundTrafficRepository InboundTrafficRepository,
	vpnTrafficService VpnTrafficService,
	userService UserTrafficService,
) *InboundTraffic {
	return &InboundTraffic{
		userService: userService,
		vpnService:  vpnTrafficService,
		inboundRepo: inboundTrafficRepository,
	}
}

func (i *InboundTraffic) UpdateTraffics() {
	active := true
	inbounds, err := i.inboundRepo.Filter(&inboundRepoDto.FilterInbound{IsActive: &active})
	if err != nil {
		return
	}

	proxy, err := i.vpnService.MakeProxy()
	if err != nil {
		return
	}

	for _, inbound := range inbounds {
		i.updateTraffic(inbound, proxy)
	}
}

func (i *InboundTraffic) updateTraffic(inbound *entity.Inbound, proxy adapter.ProxyVPN) {
	user, err := i.userService.FindByID(inbound.UserID)
	if err != nil {
		return
	}

	vpnProxyInput := adapter.GenerateVPNProxyInput(inbound, user)
	traffic, err := proxy.GetTraffic(vpnProxyInput)
	if err != nil {
		return
	}

	i.inboundRepo.IncreaseTrafficUsage(
		inbound.ID,
		uint32(traffic.Download)+uint32(traffic.Upload),
	)
}
