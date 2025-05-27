package inbound

import (
	"github.com/mohamadrezamomeni/momo/adapter"
	"github.com/mohamadrezamomeni/momo/entity"
	vpnProxy "github.com/mohamadrezamomeni/momo/proxy/vpn"
)

func (i *Inbound) HealingUpInboundExpired() {
	inbounds, err := i.inboundRepo.RetriveActiveInboundExpired()
	if err != nil {
		return
	}
	proxy, err := i.vpnService.MakeProxy()
	if err != nil {
		return
	}
	defer proxy.Close()
	for _, inbound := range inbounds {
		i.deActiveInbound(inbound, proxy)
	}
}

func (i *Inbound) HealingUpInboundOverQuoted() {
	inbounds, err := i.inboundRepo.RetriveActiveInboundsOverQuota()
	if err != nil {
		return
	}
	i.deactiveInbounds(inbounds)
}

func (i *Inbound) HealingUpInboundBlocked() {
	inbounds, err := i.inboundRepo.RetriveActiveInboundBlocked()
	if err != nil {
		return
	}
	i.deactiveInbounds(inbounds)
}

func (i *Inbound) HealingUpInboundCharged() {
	inbounds, err := i.inboundRepo.RetriveDeactiveInboundsCharged()
	if err != nil {
		return
	}
	i.activeInbounds(inbounds)
}

func (i *Inbound) activeInbounds(inbounds []*entity.Inbound) {
	proxy, err := i.vpnService.MakeProxy()
	if err != nil {
		return
	}
	defer proxy.Close()
	for _, inbound := range inbounds {
		i.activeInbound(inbound, proxy)
	}
}

func (i *Inbound) deactiveInbounds(inbounds []*entity.Inbound) {
	proxy, err := i.vpnService.MakeProxy()
	if err != nil {
		return
	}
	defer proxy.Close()
	for _, inbound := range inbounds {
		i.deActiveInbound(inbound, proxy)
	}
}

func (i *Inbound) deActiveInbound(inbound *entity.Inbound, vpnProxy vpnProxy.IProxyVPN) error {
	user, err := i.userService.FindByID(inbound.UserID)
	if err != nil {
		return err
	}

	vpnProxyInput := adapter.GenerateVPNProxyInput(inbound, user)
	err = vpnProxy.DisableInbound(vpnProxyInput)
	if err != nil {
		return err
	}

	return i.inboundRepo.DeActive(inbound.ID)
}

func (i *Inbound) activeInbound(inbound *entity.Inbound, vpnProxy vpnProxy.IProxyVPN) error {
	user, err := i.userService.FindByID(inbound.UserID)
	if err != nil {
		return err
	}

	vpnProxyInput := adapter.GenerateVPNProxyInput(inbound, user)

	if err != nil {
		return err
	}
	err = vpnProxy.AddInbound(vpnProxyInput)
	if err != nil {
		return err
	}

	return i.inboundRepo.Active(inbound.ID)
}
