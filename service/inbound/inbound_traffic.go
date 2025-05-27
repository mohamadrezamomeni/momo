package inbound

import (
	"strconv"

	"github.com/mohamadrezamomeni/momo/adapter"
	inboundRepoDto "github.com/mohamadrezamomeni/momo/dto/repository/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
)

func (i *Inbound) UpdateTraffics() {
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

func (i *Inbound) updateTraffic(inbound *entity.Inbound, proxy adapter.ProxyVPN) {
	info, err := i.getInfo(inbound)
	if err != nil {
		return
	}

	traffic, err := proxy.GetTraffic(info)
	if err != nil {
		return
	}

	i.inboundRepo.IncreaseTrafficUsage(
		strconv.Itoa(inbound.ID),
		uint32(traffic.Download)+uint32(traffic.Upload),
	)
}
